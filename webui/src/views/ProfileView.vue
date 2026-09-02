<template>
  <div class="profile-wrapper">

    <div class="profile-card">

      <h1 class="title">Profilo utente</h1>

      <!-- FOTO PROFILO -->
      <div class="profile-section">
        <h2 class="section-title">Immagine profilo</h2>

        <!-- Immagine vuota -->
        <div class="profile-photo placeholder-photo"></div>

        <button class="btn-primary">
          Cambia foto profilo
        </button>
      </div>

      <!-- USERNAME -->
      <div class="profile-section">
        <h2 class="section-title">Username</h2>

        <div class="username-box">
          <span class="placeholder-text">username</span>
        </div>

        <button class="btn-primary">
          Cambia username
        </button>
      </div>

      <!-- AZIONI -->
      <div class="actions">
        <button class="btn-secondary">
          Torna alla Home
        </button>

        <button class="btn-danger">
          Logout
        </button>
      </div>

    </div>

  </div>
</template>

<style scoped>
.profile-wrapper {
  height: 100vh;
  width: 100%;
  display: flex;
  justify-content: center;   /* CENTRA ORIZZONTALMENTE */
  align-items: center;       /* CENTRA VERTICALMENTE */
  background-color: #e2fff8; /* opzionale, come HomeView */
}

.profile-card {
  background: white;
  padding: 40px;
  border-radius: 22px;
  width: 450px;
  box-shadow: 0 20px 40px rgba(0,0,0,0.18);
  display: flex;
  flex-direction: column;
  gap: 30px;
}

.title {
  text-align: center;
  font-size: 32px;
  font-weight: bold;
}

.profile-section {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.section-title {
  font-size: 20px;
  font-weight: bold;
}

/* Foto profilo vuota */
.placeholder-photo {
  width: 140px;
  height: 140px;
  border-radius: 50%;
  background-color: #ddd;
  border: 3px solid #ccc;
}

.placeholder-text {
  color: #777;
  font-size: 20px;
}

.username-box {
  padding: 12px 18px;
  background-color: #f5f5f5;
  border-radius: 12px;
  font-size: 20px;
  font-weight: bold;
  border: 2px solid #ccc;
}

.actions {
  display: flex;
  justify-content: space-between;
  margin-top: 20px;
}

.btn-primary {
  padding: 12px 18px;
  background-color: #222;
  color: white;
  border: none;
  border-radius: 14px;
  font-size: 16px;
  font-weight: bold;
  cursor: pointer;
  transition: 0.25s ease;
}

.btn-primary:hover {
  background-color: #444;
  transform: translateY(-3px);
}

.btn-secondary {
  padding: 12px 18px;
  background-color: #ccc;
  border: none;
  border-radius: 14px;
  font-size: 16px;
  cursor: pointer;
  transition: 0.25s ease;
}

.btn-secondary:hover {
  background-color: #bbb;
}

.btn-danger {
  padding: 12px 18px;
  background-color: #d9534f;
  color: white;
  border: none;
  border-radius: 14px;
  font-size: 16px;
  cursor: pointer;
  transition: 0.25s ease;
}

.btn-danger:hover {
  background-color: #c9302c;
}
</style>
