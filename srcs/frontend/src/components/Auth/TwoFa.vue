<template>
  <div class="twoFa-container">
    <div class="twoFa-content">
      <div class="qrcode" v-if="showQrCode"></div>
      <div class="confirmTwoFA" v-if="showQrCode">
        <label for="confirmeTwoFA">Please enter the confirmation code: </label>
        <form class="container-input-confirm-2fa" @submit.prevent="confirm2FA">
          <input 
            v-model="confirmationCode" 
            type="text" 
            name="confirmeTwoFA" 
            required 
            placeholder="Enter 2FA code" 
          />
          <input type="submit" value="Send" />
        </form>
      </div>
      <div class="twoFa" v-else>
        <button @click="generateQrcode">Generate Google Authenticator QR code</button>
      </div>
      <div v-if="errorMessage" class="error-message">{{ errorMessage }}</div>
      <div v-if="successMessage" class="success-message">{{ successMessage }}</div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { API_BASE_URL } from '../../services/apiUtils';

const confirmationCode = ref(''); // Stocke le code 2FA entré par l'utilisateur
const errorMessage = ref('');
const successMessage = ref('');
const showQrCode = ref<boolean>(false); // État pour afficher ou non le QR code et le formulaire

const resetMessages = () => {
  errorMessage.value = '';
  successMessage.value = '';
};

const generateQrcode = async () => {
  resetMessages();
  // Afficher le QR code et le formulaire de confirmation
  showQrCode.value = true;
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

    // Sélectionne la div contenant le QR code
    const qrCodeDiv = document.querySelector('.qrcode');

    // Vide la div avant d'ajouter le nouveau QR code
    if (qrCodeDiv) {
      qrCodeDiv.innerHTML = '';
    }

    // Crée une balise <img> pour afficher le nouveau QR code
    const imgElement = document.createElement('img');
    imgElement.src = qrCodeUrl;
    imgElement.alt = 'Scan this QR Code with Google Authenticator';

    // Ajoute l'image à la div
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
    successMessage.value = `2FA Confirmation successful: ${result.message}`;
  } catch (error) {
    console.error('Error confirming 2FA:', error);
    errorMessage.value = "Failed to confirm 2FA. Please try again.";
  }
};
</script>

<style scoped>
.twoFa-container {
  font-weight: 400;
  font-style: normal;
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
  min-height: 250px;
}
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
  position: relative;
  bottom: 100%;
  left: 0;
  right: 0;
  padding: 8px;
  border-radius: 4px;
  font-size: 14px;
  text-align: center;
  z-index: 1002;
  margin-bottom: 8px;
  max-height: none;
  overflow-wrap: break-word;
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
</style>
