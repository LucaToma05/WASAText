<template>
  <div class="overlay">
    <div class="overlay-box">
      <h3>Crea nuovo gruppo</h3>

      <input 
        type="text"
        v-model="groupName"
        placeholder="Nome del gruppo"
      />

      <input 
        type="text"
        v-model="groupUsers"
        placeholder="Username utenti (separati da virgola)"
      />

      <div class="overlay-actions">
        <button @click="createGroup">Crea</button>
        <button @click="$emit('close')">Annulla</button>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      groupName: "",
      groupUsers: ""
    };
  },

  methods: {
    async createGroup() {
      if (!this.groupName.trim()) return;

      const users = this.groupUsers
        .split(",")
        .map(u => u.trim())
        .filter(u => u.length > 0);

      try {
        await this.$axios.post("/createGroup", {
          name: this.groupName,
          users: users
        });

        this.$emit("created");
        this.$emit("close");
      } catch (err) {
        console.log("Errore creazione gruppo:", err);
      }
    }
  }
};
</script>

<style scoped>
.overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0,0,0,0.4);
  display: flex;
  justify-content: center;
  align-items: center;
}

.overlay-box {
  background: white;
  padding: 25px;
  border-radius: 8px;
  width: 350px;
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.overlay-actions {
  display: flex;
  justify-content: space-between;
}
</style>
