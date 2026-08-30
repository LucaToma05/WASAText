<script>
export default {
  data() {
    return {
      userId: null,
      userData: null,
      loading: true,
      errorMsg: ""
    };
  },

  async mounted() {
    this.userId = localStorage.getItem("userId");
    await this.loadProfile();
  },

  methods: {
    async loadProfile() {
      try {
        const response = await this.$axios.get(`/getUser/${this.userId}`, {
          headers: { Authorization: "Bearer " + this.userId }
        });

        this.userData = response.data;
        this.loading = false;
      } catch (err) {
        console.log("Errore caricamento profilo:", err);
        this.errorMsg = "Impossibile caricare il profilo";
        this.loading = false;
      }
    },

    logout() {
      localStorage.removeItem("userId");
      this.$router.push("/");
    }
  }
};
</script>

<template>
  <div class="profile-wrapper">

    <div class="profile-card" v-if="!loading && userData">
      <h1>Profilo</h1>

      <img 
        class="profile-photo"
        :src="userData.photoUrl || 'https://via.placeholder.com/120'"
      />

      <p class="profile-name">{{ userData.name }}</p>
      <p class="profile-id">ID: {{ userData.id }}</p>

      <div class="profile-actions">
        <button @click="logout">Logout</button>
        <button @click="$router.push('/home')">Torna alla Home</button>
      </div>
    </div>

    <div v-if="loading" class="loading">
      Caricamento profilo...
    </div>

    <div v-if="errorMsg" class="error">
      {{ errorMsg }}
    </div>

  </div>
</template>

<style scoped>
.profile-wrapper {
  height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
}

.profile-card {
  background: white;
  padding: 40px;
  border-radius: 18px;
  box-shadow: 0 20px 40px rgba(0,0,0,0.18);
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 25px;
  width: 420px;
}

.profile-photo {
  width: 120px;
  height: 120px;
  border-radius: 50%;
}

.profile-name {
  font-size: 28px;
  font-weight: bold;
}

.profile-id {
  font-size: 16px;
  color: #666;
}

.profile-actions {
  display: flex;
  gap: 20px;
}

.profile-actions button {
  padding: 10px 18px;
  border: none;
  background-color: #222;
  color: white;
  font-weight: bold;
  cursor: pointer;
  border-radius: 8px;
}

.profile-actions button:hover {
  background-color: #444;
}

.loading {
  font-size: 22px;
  color: #333;
}

.error {
  font-size: 20px;
  color: red;
}
</style>
