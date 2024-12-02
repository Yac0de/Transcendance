<template>
    <div class="twoFa-container">
      <div class="twoFa-content">
            <div class="qrcode">
            </div>
            <div class="confirmTwoFA">
                <label for="confirmeTwoFA">Please enter the confirmation code: </label>
                <div class="container-input-confirm-2fa">
                    <input type="text" name="confirmeTwoFA" required />
                    <input type="submit" value="send"/>
                </div>
                <div class="twoFa">
                    <button @click = "generateQrcode">Generate Google Authenticator QR code</button>
                </div>
            </div>
      </div>
    </div>
</template>

<script setup lang="ts">

import { API_BASE_URL } from '../../services/apiUtils';

const generateQrcode = async () => {
    fetch(`${API_BASE_URL}/auth/generate2FA`, {
        method: 'GET',
        credentials: 'include',
    })
    .then(response => {
        if (!response.ok) {
            throw new Error('Failed to fetch 2FA QR code');
        }
        return response.blob();
    })
    .then(blob => {
        // Créer une URL temporaire pour le Blob
        const qrCodeUrl = URL.createObjectURL(blob);

        // Afficher l'image dans une balise <img>
        const imgElement = document.createElement('img');
        imgElement.src = qrCodeUrl;
        imgElement.alt = 'Scan this QR Code with Google Authenticator';
        const div = document.querySelector(".qrcode");
        div?.appendChild(imgElement);
    })
    .catch(error => {
        console.error('Error fetching 2FA QR code:', error);
    });
}
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
    justify-content: end;
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
    background: linear-gradient(to right, var(--secondary-bright-color), color-mix(in srgb, var(--secondary-bright-color) 75%, white));
    color: white;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    font-size: 16px;
    margin-top: 20px;
  }

  .twoFa-content button:hover {
    background: linear-gradient(to right, var(--secondary-bright-color), color-mix(in srgb, var(--secondary-bright-color) 85%, white));
    transform: scale(1.02);
  }
  
  .confirmTwoFA {
    display: flex;
    flex-direction: column;
    color: white;
  }

  .container-input-confirm-2fa {
    display: flex;
  }
  </style>
  
  
  