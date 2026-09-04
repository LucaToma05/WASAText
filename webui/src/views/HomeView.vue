<template>
  <div class="home-layout">

    <!-- SIDEBAR -->
    <div class="home-sidebar">

      <!-- Barra dei pulsanti -->
      <div class="sidebar-top-bar">
        <div class="left-buttons">
          <button class="sidebar-btn" @click="showCreateChat = true"> Nuova chat </button>
          <button class="sidebar-btn">Nuovo gruppo</button>
        </div>

        <div class="right-button">
          <button class="sidebar-btn" @click="goToProfile">Profilo</button>
        </div>
      </div>

      <!-- LISTA CHAT -->
      <div class="sidebar-content">
        <div 
          v-for="chat in sortedChats" 
          :key="chat.chatId" 
          class="chat-preview"
          :class="{ selected: chat.chatId === selectedChatId }"
          @click="openChat(chat.chatId)"
        >
          <img :src="chat.photo || '/default.png'" class="chat-photo">

          <div class="chat-info">
            <div class="chat-top-row">
              <span class="chat-name">{{ chat.name }}</span>
              <span class="chat-time">{{ formatTime(chat.lastMessageTime) }}</span>
            </div>

            <div class="chat-last-message">
              {{ chat.lastMessage }}
            </div>
          </div>
        </div>
      </div>

    </div>

    <!-- AREA DESTRA -->
    <div class="home-main-area">

      <!-- Se nessuna chat è selezionata -->
      <div v-if="!selectedChat" class="empty-chat">
        <p>Seleziona una chat dalla sidebar.</p>
      </div>

      <!-- CHAT APERTA -->
      <div v-else class="chat-area">

        <!-- HEADER -->
        <div class="chat-header">
          <h2>{{ selectedChat.name }}</h2>
        </div>

        <!-- MESSAGGI -->
        <div class="messages">
          <div 
            v-for="msg in selectedChat.messages"
            :key="msg.id"
            class="message"
            :class="{ mine: msg.sender === userId }"
          >
            <p>{{ msg.text }}</p>
            <span class="timestamp">{{ formatTime(msg.instant) }}</span>
          </div>

          <div v-if="selectedChat.messages.length === 0" class="no-messages">
            Nessun messaggio ancora. Inizia la conversazione!
          </div>
        </div>

        <!-- BOX INVIO -->
        <div class="send-box">
          <input 
            type="text"
            v-model="newMessage"
            placeholder="Scrivi un messaggio..."
            @keyup.enter="sendMessage"
          />
          <button @click="sendMessage">Invia</button>
        </div>

      </div>

    </div>

    <!-- MODAL CREA CHAT -->
    <div v-if="showCreateChat" class="modal-overlay" @click.self="closeCreateChat">
      <div class="modal-card">

        <h2 class="modal-title">Crea nuova chat</h2>

        <input 
          v-model="newChatUsername"
          type="text"
          placeholder="Username dell'utente"
          class="modal-input"
        />

        <p v-if="errorMsg" class="error-msg">{{ errorMsg }}</p>

        <div class="modal-buttons">
          <button class="btn-primary" @click="createChat">Crea chat</button>
          <button class="btn-secondary" @click="closeCreateChat">Annulla</button>
        </div>

      </div>
    </div>

  </div>
</template>

<script>
export default {
  data() {
    return {
      showCreateChat: false,
      newChatUsername: "",
      errorMsg: "",
      chats: [],
      selectedChat: null,
      selectedChatId: null,
      newMessage: "",
      userId: localStorage.getItem("userId")
    };
  },

  computed: {
    sortedChats() {
      return [...this.chats].sort((a, b) => {
        return new Date(b.lastMessageTime || 0) - new Date(a.lastMessageTime || 0);
      });
    }
  },

  methods: {
    goToProfile() {
      this.$router.push("/profile");
    },

    async openChat(chatId) {
      this.selectedChatId = chatId;

      try {
        const response = await this.$axios.get(`/chats/${chatId}/messages`, {
          headers: {
            Authorization: "Bearer " + localStorage.getItem("token")
          }
        });

        const messages = Array.isArray(response.data) ? response.data : [];

        this.selectedChat = {
          chatId,
          name: this.chats.find(c => c.chatId === chatId)?.name || "Chat",
          messages
        };

      } catch (err) {
        console.log("Errore apertura chat:", err);
      }
    },

    formatTime(timeString) {
      if (!timeString) return "";
      const date = new Date(timeString);
      return date.toLocaleTimeString("it-IT", {
        hour: "2-digit",
        minute: "2-digit"
      });
    },

    async loadChats() {
      try {
        const response = await this.$axios.get(`/chats`, {
          headers: {
            Authorization: "Bearer " + localStorage.getItem("token")
          }
        });

        this.chats = response.data.map(chat => ({
          ...chat,
          displayName: chat.name
        }));

      } catch (err) {
        console.log("Errore caricamento chat:", err);
      }
    },

    async createChat() {
      if (this.newChatUsername.trim() === "") {
        this.errorMsg = "Inserisci un username";
        return;
      }

      try {
        const response = await this.$axios.post(
          "/chats",
          {
            type: false,
            name: this.newChatUsername
          },
          {
            headers: {
              Authorization: "Bearer " + localStorage.getItem("token")
            }
          }
        );

        this.chats.unshift(response.data);
        this.closeCreateChat();

      } catch (err) {
        this.errorMsg = err.response?.data || "Errore creazione chat";
      }
    },

    async sendMessage() {
      if (this.newMessage.trim() === "") return;

      try {
        await this.$axios.post(
          `/chats/${this.selectedChatId}/messages`,
          { text: this.newMessage },
          {
            headers: {
              Authorization: "Bearer " + localStorage.getItem("token")
            }
          }
        );

        this.newMessage = "";
        await this.openChat(this.selectedChatId);

      } catch (err) {
        console.log("Errore invio messaggio:", err);
      }
    },

    closeCreateChat() {
      this.showCreateChat = false;
      this.newChatUsername = "";
      this.errorMsg = "";
    }
  },

  mounted() {
    this.loadChats();

    this.intervalId = setInterval(() => {
      this.loadChats();
    }, 3000);
  },

  beforeUnmount() {
    clearInterval(this.intervalId);
  }
};
</script>

<style scoped>
/* Layout orizzontale */
.home-layout {
  display: flex;
  width: 100%;
  height: 100%;
  background-color: #e2fff8;
}

/* Sidebar */
.home-sidebar {
  width: 30%;
  background-color: #ffffff;
  border-right: 2px solid #ccc;
  display: flex;
  flex-direction: column;

  box-shadow: 0 20px 40px rgba(0,0,0,0.18);
  border-radius: 0 22px 22px 0;
}

/* Barra superiore */
.sidebar-top-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;

  padding: 20px;
  border-bottom: 2px solid #ccc;

  background-color: #ffffff;
  border-radius: 0 22px 0 0;
}

.left-buttons {
  display: flex;
  gap: 12px;
}

.right-button {
  display: flex;
}

/* Bottoni */
.sidebar-btn {
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

.sidebar-btn:hover {
  background-color: #444;
  transform: translateY(-3px);
  box-shadow: 0 10px 20px rgba(0,0,0,0.25);
}

/* LISTA CHAT */
.sidebar-content {
  flex: 1;
  overflow-y: auto;
  padding: 10px 0;
}

.chat-preview {
  display: flex;
  gap: 14px;
  padding: 16px;
  cursor: pointer;
  transition: 0.25s ease;
  border-bottom: 1px solid #eee;
}

.chat-preview.selected {
  background-color: #d0f0ff;
  border-left: 5px solid #0077aa;
}


.chat-preview:hover {
  background-color: #f5f5f5;
}

.chat-photo {
  width: 55px;
  height: 55px;
  border-radius: 50%;
}

.chat-info {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.chat-top-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.chat-name {
  font-weight: bold;
  font-size: 18px;
  color: #222;
}

.chat-time {
  font-size: 14px;
  color: #666;
}

.chat-last-message {
  font-size: 15px;
  color: #444;
  margin-top: 4px;
}

/* Area principale */
.home-main-area {
  width: 70%;
  background-color: #d9f7ff;
  border-radius: 22px;
  margin-left: 20px;

  box-shadow: 0 20px 40px rgba(0,0,0,0.18);
}

.empty-chat {
  height: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
  text-align: center;

  font-size: 22px;
  color: #444;
  padding: 20px;
}


/* MODAL */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0,0,0,0.45);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 500;
}

.modal-card {
  background-color: #ffffff;
  padding: 40px;
  border-radius: 22px;
  width: 420px;

  box-shadow: 0 20px 40px rgba(0,0,0,0.18);
  display: flex;
  flex-direction: column;
  gap: 25px;
}

.modal-title {
  font-size: 32px;
  font-weight: bold;
  color: #222;
  text-align: center;
}

.modal-input {
  width: 100%;
  height: 60px;
  font-size: 22px;
  padding: 12px;
  border-radius: 14px;
  border: 2px solid #555;
  background-color: #f5f5f5;
  transition: 0.25s ease;
}

.modal-input:focus {
  border-color: #000;
  background-color: #fff;
  outline: none;
  box-shadow: 0 0 10px rgba(0,0,0,0.15);
}

.error-msg {
  color: red;
  font-size: 18px;
  text-align: center;
}

.modal-buttons {
  display: flex;
  justify-content: space-between;
}

.btn-primary {
  padding: 12px 18px;
  background-color: #222;
  color: white;
  border: none;
  border-radius: 14px;
  font-size: 18px;
  font-weight: bold;
  cursor: pointer;
  transition: 0.25s ease;
}

.btn-primary:hover {
  background-color: #444;
  transform: translateY(-3px);
  box-shadow: 0 10px 20px rgba(0,0,0,0.25);
}

.btn-secondary {
  padding: 12px 18px;
  background-color: #ccc;
  border: none;
  border-radius: 14px;
  font-size: 18px;
  cursor: pointer;
  transition: 0.25s ease;
}

.btn-secondary:hover {
  background-color: #bbb;
}

/* AREA CHAT */
.chat-area {
  height: 100%;
  display: flex;
  flex-direction: column;
  background-color: #ffffff;
  border-radius: 22px;
}

/* Header */
.chat-header {
  padding: 20px;
  border-bottom: 2px solid #ccc;
  font-size: 22px;
  font-weight: bold;
}

/* Messaggi */
.messages {
  flex: 1;
  overflow-y: auto;
  padding: 20px;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.message {
  max-width: 60%;
  padding: 12px;
  background-color: #e6e6e6;
  border-radius: 12px;
}

.message.mine {
  background-color: #cce5ff;
  margin-left: auto;
}

.timestamp {
  font-size: 12px;
  color: #666;
  margin-top: 4px;
}

/* Box invio */
.send-box {
  display: flex;
  padding: 15px;
  border-top: 2px solid #ccc;
  background-color: #f5f5f5;
}

.send-box input {
  flex: 1;
  padding: 12px;
  border-radius: 12px;
  border: 1px solid #aaa;
}

.send-box button {
  margin-left: 10px;
  padding: 12px 20px;
  border-radius: 12px;
  background-color: #222;
  color: white;
  font-weight: bold;
  cursor: pointer;
}



/* MOBILE */
@media (max-width: 900px) {
  .home-layout {
    flex-direction: column;
    padding: 10px;
  }

  .home-sidebar {
    width: 100%;
    border-radius: 22px;
    border-right: none;
    border-bottom: 2px solid #ccc;
  }

  .sidebar-top-bar {
    flex-direction: column;
    gap: 12px;
    align-items: stretch;
  }

  .left-buttons,
  .right-button {
    justify-content: center;
  }

  .sidebar-btn {
    width: 100%;
    text-align: center;
  }

  .home-main-area {
    width: 100%;
    margin-left: 0;
    margin-top: 20px;
    border-radius: 22px;
  }

  .modal-card {
    width: 90%;
  }
}
</style>