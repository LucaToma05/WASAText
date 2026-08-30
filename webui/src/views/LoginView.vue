<template>
  <div class="login-page">
    <div class="login-container">
      <h1 id="login-label">LOGIN</h1>
      <form @submit.prevent="login">
        <input v-model="username" type="text" placeholder="Username">
        <button type="submit">Accedi</button>
      </form>
      <p v-if="errorMsg">{{ errorMsg }}</p>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      username: "",
      errorMsg: ""
    };
  },

  methods: {
    async login() {
      if (this.username.trim() === "") {
        this.errorMsg = "Username vuoto";
        return;
      }

      try {
        const response = await this.$axios.post("/session", { 
          name: this.username
        });

        const userId = response.data.id;
        localStorage.setItem("userId", userId);
        this.$router.push("/home");
      } catch (err) {
        this.errorMsg = err.response?.data?.error || "Errore durante il login";
      }
    }
  }
};
</script>

<style scoped>
/* riempie il container sotto l'header */
.login-page {
  width: 100%;
  height: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
}

/* Card del login */
.login-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 40px;

  padding: 70px;
  border-radius: 22px;
  background-color: #ffffff;

  box-shadow: 0 20px 40px rgba(0,0,0,0.18);
}

/* Titolo */
#login-label {
  font-size: 72px;
  letter-spacing: 4px;
  font-weight: bold;
  color: #222;
}

.login-container form {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.login-container input,
.login-container button {
  width: 420px;
  height: 75px;
  font-size: 28px;
  padding: 14px;
  border-radius: 14px;
}

.login-container input {
  border: 2px solid #555;
  background-color: #f5f5f5;
  transition: 0.25s ease;
}

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
  cursor: pointer;
  transition: 0.25s ease;
}

.login-container button:hover {
  background-color: #444;
  transform: translateY(-3px);
  box-shadow: 0 10px 20px rgba(0,0,0,0.25);
}
</style>
