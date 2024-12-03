<template>
  <div class="twoFa-content" v-if="!isLoading">
    <div v-if="showQrCode && !is2FAActive" class="qrcode"></div>
    <div v-if="showQrCode && !is2FAActive" class="confirmTwoFA">
      <label for="confirmationCode">{{ $t('enterConfirmationCode') }}</label>
      <form
        class="container-input-confirm-2fa"
        @submit.prevent="confirm2FA"
      >
        <input
          v-model="confirmationCode"
          type="text"
          id="confirmationCode"
          required
          :placeholder="$t('enter2FACodePlaceholder')"
        />
        <input type="submit" :value="$t('sendButton')" />
      </form>
    </div>
    <div v-if="!showQrCode && !is2FAActive" class="twoFa">
      <button @click="generateQrcode">
        {{ $t('generateQRCode') }}
      </button>
    </div>
    <div
      v-if="errorMessage || successMessage"
      :class="{
        'error-message': errorMessage,
        'success-message': successMessage
      }"
    >
      {{ errorMessage || successMessage }}
    </div>
  </div>
  <div v-else class="loading">
    {{ $t('loading') }}
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useI18n } from 'vue-i18n';
import { API_BASE_URL } from '../../services/apiUtils';

const confirmationCode = ref('');
const errorMessage = ref('');
const successMessage = ref('');
const isLoading = ref(true);
const showQrCode = ref<boolean>(false);
const is2FAActive = ref<boolean>(false);

const { t } = useI18n();

const resetMessages = () => {
  errorMessage.value = '';
  successMessage.value = '';
};

const check2FaStatus = async () => {
  resetMessages();
  try {
    const response = await fetch(`${API_BASE_URL}/auth/2FA-status`, {
      method: 'GET',
      credentials: 'include',
    });

    if (!response.ok) {
      throw new Error('Failed to fetch 2FA status');
    }

    const result = await response.json();
    is2FAActive.value = result.status;

    if (is2FAActive.value) {
      errorMessage.value = t('alreadyActive2FA');
    }

  } catch (error) {
    console.error('Error fetching 2FA status:', error);
    errorMessage.value = t('fetchStatusError');
  }
  finally {
    isLoading.value = false;
  }
};

const generateQrcode = async () => {
  resetMessages();
  showQrCode.value = true;

  if (is2FAActive.value) {
    errorMessage.value = t('alreadyActive2FA');
    return;
  }

  try {
    const response = await fetch(`${API_BASE_URL}/auth/generate2FA`, {
      method: 'GET',
      credentials: 'include',
    });

    if (!response.ok) {
      throw new Error('Failed to fetch 2FA QR code');
    }

    const blob = await response.blob();
    const qrCodeUrl = URL.createObjectURL(blob);

    const qrCodeDiv = document.querySelector('.qrcode');
    if (qrCodeDiv) {
      qrCodeDiv.innerHTML = '';
    }

    const imgElement = document.createElement('img');
    imgElement.src = qrCodeUrl;
    imgElement.alt = 'Scan this QR Code with Google Authenticator';
    qrCodeDiv?.appendChild(imgElement);

  } catch (error) {
    console.error('Error fetching 2FA QR code:', error);
    errorMessage.value = 'Failed to generate 2FA QR code. Please try again.';
  }
};

const confirm2FA = async () => {
  resetMessages();
  try {
    const response = await fetch(`${API_BASE_URL}/auth/verify2FA`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify({ code: confirmationCode.value }),
    });

    if (!response.ok) {
      throw new Error('Failed to confirm 2FA code');
    }

    const result = await response.json();
    successMessage.value = t('confirm2FASuccess', { message: result.message });
    is2FAActive.value = true;
  } catch (error) {
    console.error('Error confirming 2FA:', error);
    errorMessage.value = t('confirm2FAError');
  }
};

onMounted(() => {
  check2FaStatus();
});
</script>

<style scoped>
.twoFa-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: space-around;
  text-align: center;
  height: 60vh;
  min-height: 65px;
  min-width: 300px;
  padding: 2vh 5vw;
  border-radius: 20px;
  box-shadow: 0 0 20px rgba(0, 0, 0, 0.5);
  overflow: hidden;
  background: var(--main-color);
}
.twoFa-content button {
  width: 100%;
  padding: 10px;
  background: linear-gradient(
    to right,
    var(--secondary-bright-color),
    color-mix(in srgb, var(--secondary-bright-color) 75%, white)
  );
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 16px;
  margin-top: 20px;
}
.twoFa-content button:hover {
  background: linear-gradient(
    to right,
    var(--secondary-bright-color),
    color-mix(in srgb, var(--secondary-bright-color) 85%, white)
  );
  transform: scale(1.02);
}
.confirmTwoFA {
  display: flex;
  flex-direction: column;
  color: white;
}
.container-input-confirm-2fa {
  display: flex;
  justify-content: space-between;
}
.container-input-confirm-2fa input[type='text'] {
  width: 100%;
}
.container-input-confirm-2fa input[type='submit'] {
  width: 25%;
  padding: 10px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  color: white;
  font-size: 14px;
  transition: background-color 0.3s;
  background: linear-gradient(
    to right,
    var(--secondary-dark-color),
    color-mix(in srgb, var(--secondary-dark-color) 75%, white)
  );
}

.error-message,
.success-message {
  position: static;
  margin: 16px 0;
  padding: 8px;
  border-radius: 4px;
  font-size: 14px;
  text-align: center;
  word-wrap: break-word;
}


.error-message {
  background-color: #ffebee;
  color: #d32f2f;
  border: 1px solid #ef9a9a;
}

.success-message {
  background-color: #e8f5e9;
  color: #388e3c;
  border: 1px solid #a5d6a7;
}

.loading {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100vh;
  font-size: 18px;
  font-weight: bold;
  color: var(--secondary-dark-color);
}

</style>
