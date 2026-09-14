<template>
  <div class="login-page">
    <div class="login-container">
      <h1 id="login-label">LOGIN</h1>
      <!-- Form che chiama la funzione login quando viene inviato e .prevent impedisce il reload della pagina -->
      <form @submit.prevent="login">
        <!-- Input username con legame bidirezionale -->
        <input v-model="username" type="text" placeholder="Username">
        <!-- Bottone che chiama la funzione login quando viene cliccato -->
        <button type="submit">Accedi</button>
      </form>

      <!-- Messaggio di errore, mostrato solo se presente e non "" -->
      <p v-if="errorMsg">{{ errorMsg }}</p>
    </div>
  </div>
</template>

<script>
export default {
  // Dati reattivi
  data() {
    return {
      username: "",
      errorMsg: ""
    };
  },

  methods: {
    async login() {
      // Se l'username è vuoto (o contiene solo spazi), mostra un errore e interrompe la funzione
      if (this.username.trim() === "") {
        this.errorMsg = "Username vuoto";
        return;
      }

      // Tentativo di login
      try {
        const response = await this.$axios.post("/session", { 
          name: this.username
        });

        // Salva i dati in maniera persistente
        localStorage.setItem("token", response.data.token);
        localStorage.setItem("username", this.username.trim());

        // L'utente loggato viene reindirizzato alla home
        this.$router.push("/home");
      } catch (err) {
        this.errorMsg = err.response?.data?.error || "Errore durante il login";
      }
    }
  }
};
</script>

<style scoped>
.login-page {
  width: 100%;
  height: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 16px;
}

.login-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 30px;

  padding: clamp(30px, 6vw, 70px);
  border-radius: 22px;
  background-color: #ffffff;
  width: 100%;
  max-width: 520px;

  box-shadow: 0 20px 40px rgba(0,0,0,0.18);
}

#login-label {
  font-size: clamp(36px, 10vw, 72px);
  letter-spacing: 4px;
  font-weight: bold;
  color: #222;
}

.login-container form {
  display: flex;
  flex-direction: column;
  gap: 20px;
  width: 100%;
}

.login-container input,
.login-container button {
  width: 100%;
  max-width: 420px;
  height: clamp(50px, 8vw, 75px);
  font-size: clamp(16px, 4vw, 28px);
  padding: 14px;
  border-radius: 14px;
  align-self: center;
}

.login-container input {
  border: 2px solid #555;
  background-color: #f5f5f5;
  transition: 0.25s ease;
}

/* Quando viene cliccato */ 
.login-container input:focus {
  border-color: #000;
  background-color: #fff;
  outline: none;
  box-shadow: 0 0 10px rgba(0,0,0,0.15);
}

.login-container button {
  border: none;
  background-color: #222;
  color: white;
  font-weight: bold;
  transition: 0.25s ease;
}

/* Quando il cursore passa */
.login-container button:hover {
  background-color: #444;
  transform: translateY(-3px);
  box-shadow: 0 10px 20px rgba(0,0,0,0.25);
}

/* Quando la schermata è troppo piccola */ 
@media (max-width: 600px) {
  .login-container {
    gap: 20px;
    border-radius: 16px;
  }

  .login-container input,
  .login-container button {
    border-radius: 12px;
  }
}
</style>
