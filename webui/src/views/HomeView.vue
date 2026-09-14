<template>
  <div class="home-layout" :class="{ 'chat-open': selectedChat }">

    <!-- SIDEBAR -->
    <div class="home-sidebar">

      <!-- Barra dei pulsanti -->
      <div class="sidebar-top-bar">
        <div class="left-buttons">
          <!-- Se cliccato apre il pannello per creare una chat -->
          <button class="sidebar-btn" @click="showCreateChat = true">Nuova chat</button>
          <!-- Se cliccato apre il pannello per creare un gruppo -->
          <button class="sidebar-btn" @click="showCreateGroup = true">Nuovo gruppo</button>
        </div>

        <div class="right-button">
          <!-- Se cliccato apre la pagina di gestione del profilo-->
          <button class="sidebar-btn" @click="$router.push('/profile')">Profilo</button>
        </div>
      </div>

      <!-- LISTA CHAT da ChatItem.vue a cui passiamo proprietà e definiamo click -->
      <div class="sidebar-content">
        <ChatItem
          v-for="chat in sortedChats"
          :key="chat.chatId"
          :chat="chat"
          :is-selected="chat.chatId === selectedChatId"
          :is-unread="chat.hasUnread"
          @click="openChat(chat.chatId)"
        />
      </div>

    </div>

    <!-- AREA DESTRA -->
    <div class="home-main-area">

      <!-- Se nessuna chat è selezionata -->
      <div v-if="!selectedChat" class="empty-chat">
        <p>Seleziona una chat dalla sidebar.</p>
      </div>

      <!-- Se una chat è selezionata -->
      <div v-else class="chat-area">

        <!-- HEADER -->
        <div class="chat-header">
          <div class="chat-header-left">
            <button class="btn-back-mobile" @click="goBackToList" title="Torna alla lista chat">←</button> <!-- Solo per mobile-->
            <img 
              :src="selectedChat.photo || '/default.png'" 
              class="chat-header-photo" 
              alt="Foto conversazione"
            />
            <div class="chat-header-details">
              <h2 class="chat-header-title">{{ selectedChat.name }}</h2>
              <span v-if="selectedChat.isGroup" class="chat-header-badge">Gruppo</span>
            </div>
          </div>

          <!-- Pulsante impostazioni gruppo: visibile esclusivamente per gruppi -->
          <div v-if="selectedChat.isGroup" class="chat-header-actions">
            <button 
              type="button" 
              class="header-settings-btn" 
              @click="showGroupSettings = true" 
              title="Gestione gruppo"
            >⚙</button>
          </div>
        </div>

        <!-- MESSAGGI -->
        <div class="messages" ref="messagesContainer" @click="closeEmojiPicker">
          <!-- Itera i messaggi in chat e assegna le classi -->
          <div 
            v-for="msg in selectedChat.messages"
            :key="msg.id"
            class="message-wrapper"
            :class="{ 
              'sent': isMyMessage(msg), 
              'received': !isMyMessage(msg),
              'is-deleted': isMessageDeleted(msg),
              'has-active-picker': activeEmojiPickerMessageId === msg.id 
            }"
          >
            <!-- Messaggio -->
            <div class="message">
              <!-- Se il messaggio è stato eliminato -->
              <template v-if="isMessageDeleted(msg)">
                <div class="message-text">
                  <span class="tombstone-text"><em>Questo messaggio è stato eliminato</em></span>
                </div>
                <div class="message-meta">
                  <span class="timestamp">{{ formatTime(msg.instant) }}</span>
                </div>
              </template>

              <!-- Messaggio attivo normale -->
              <template v-else>
                <!-- Indicatore Inoltrato -->
                <div v-if="msg.forwarded" class="forwarded-indicator">
                  ↪ Inoltrato
                </div>

                <!-- Mittente: visibile solo nei messaggi ricevuti -->
                <div v-if="!isMyMessage(msg)" class="message-sender">
                  {{ getSenderName(msg) }}
                </div>

                <!-- Messaggio con risposta incorporata -->
                <div 
                  v-if="getReplyInfo(msg)" 
                  class="message-reply-preview"
                >
                  <div class="reply-preview-stripe"></div>
                  <div class="reply-preview-body">
                    <span class="reply-preview-sender">{{ getReplyInfo(msg).sender }}</span>
                    <div class="reply-preview-text">
                      <span v-if="getReplyInfo(msg).isPhoto" class="photo-indicator">📷 Foto</span>
                      <span v-else>{{ getReplyInfo(msg).text }}</span>
                    </div>
                  </div>
                </div>

                <!-- Testo del messaggio / Foto -->
                <div class="message-text">
                  <div v-if="isMessagePhoto(msg.text)" class="message-photo-wrapper">
                    <img 
                      :src="msg.text" 
                      class="message-photo-img" 
                      alt="Foto messaggio" 
                      loading="lazy"
                    />
                  </div>
                  <span v-else>{{ msg.text }}</span>
                </div>

                <!-- Reazioni -->
                <div v-if="msg.comments && msg.comments.length > 0" class="message-reactions-footer">
                  <button 
                    v-for="comment in msg.comments" 
                    :key="comment.id"
                    type="button"
                    class="reaction-footer-badge"
                    :class="{ 'is-mine': isMyReactionComment(comment) }"
                    :title="isMyReactionComment(comment) ? 'Clicca per rimuovere la tua reazione (' + comment.emoji + ')' : 'Reazione aggiunta da ' + getReactionSenderName(comment)"
                    @click.stop="handleCommentClick(msg, comment)"
                  >
                    <span class="reaction-badge-emoji">{{ comment.emoji }}</span>
                    <span class="reaction-badge-user">{{ getReactionSenderName(comment) }}</span>
                  </button>
                </div>

                <!-- Orario e spunte -->
                <div class="message-meta">
                  <span class="timestamp">{{ formatTime(msg.instant) }}</span>
                  <span 
                    v-if="isMyMessage(msg)" 
                    class="message-status" 
                    :class="getStatusClass(msg)"
                    :title="getStatusTitle(msg)"
                  >{{ isMessageRead(msg) ? '✓✓' : '✓' }}</span>
                </div>

                <!-- Barra azioni messaggio al passaggio del mouse -->
                <div class="message-action-toolbar">
                  <button type="button" class="msg-action-btn" title="Reagisci" @click.stop="toggleEmojiPicker(msg.id)">😊</button>
                  <button type="button" class="msg-action-btn" title="Rispondi" @click.stop="startReply(msg)">↩</button>
                  <button type="button" class="msg-action-btn" title="Inoltra" @click.stop="openForwardModal(msg)">↪</button>
                  <button v-if="isMyMessage(msg)" type="button" class="msg-action-btn btn-delete-msg" title="Elimina messaggio" @click.stop.prevent="deleteMessage(msg)">🗑</button>
                </div>

                <!-- Selettore emoji -->
                <div 
                  v-if="activeEmojiPickerMessageId === msg.id" 
                  class="emoji-picker-popup"
                  :class="{ 'placement-bottom': emojiPickerPlacement === 'bottom' }"
                  @click.stop
                >
                  <button 
                    v-for="emoji in ['👍', '❤️', '😂', '😮']" 
                    :key="emoji"
                    type="button"
                    class="emoji-option-btn"
                    :class="{ 'is-selected': hasUserReactedWith(msg, emoji) }"
                    @click="handleEmojiSelect(msg, emoji)"
                    :title="hasUserReactedWith(msg, emoji) ? 'Rimuovi reazione' : 'Reagisci con ' + emoji"
                  >{{ emoji }}</button>
                </div>
              </template>
            </div>

          </div>

          <!-- Messaggio se la chat non ha ancora messaggi -->
          <div v-if="selectedChat.messages.length === 0" class="no-messages">
            Nessun messaggio ancora. Inizia la conversazione!
          </div>
        </div>

        <!-- Barra preview risposta (quando si risponde a un messaggio) -->
        <div v-if="replyingToMessage" class="reply-preview-bar">
          <div class="reply-preview-stripe"></div>
          <div class="reply-preview-body">
            <div class="reply-preview-title">
              Risposta a <strong>{{ getSenderName(replyingToMessage) }}</strong>
            </div>
            <div class="reply-preview-snippet">
              <span v-if="isMessagePhoto(replyingToMessage.text)" class="photo-indicator">
                📷 Foto
              </span>
              <span v-else>{{ replyingToMessage.text }}</span>
            </div>
          </div>
          <button type="button" class="btn-cancel-reply" @click="cancelReply" title="Annulla risposta">
            ✕
          </button>
        </div>

        <!-- BOX INVIO -->
        <div class="send-box">
          <!-- Pulsante per inviare foto -->
          <button 
            type="button" 
            class="btn-attach-photo" 
            @click="showPhotoModal = true" 
            title="Invia una foto"
          >📷</button>

          <input 
            ref="messageInput"
            type="text"
            v-model="newMessage"
            placeholder="Scrivi un messaggio..."
            @keyup.enter="sendMessage"
          />
          <button class="btn-send-message" @click="sendMessage">Invia</button>
        </div>

      </div>

    </div>

    <!-- Pop-up a tutto schermo -->
    <CreateChatModal
      v-if="showCreateChat"
      @close="showCreateChat = false"
      @created="onChatCreated"
    />

    <CreateGroupModal
      v-if="showCreateGroup"
      @close="showCreateGroup = false"
      @created="onChatCreated"
    />

    <GroupSettingsModal
      v-if="showGroupSettings"
      :chat="selectedChat"
      :current-user-id="userId"
      @close="showGroupSettings = false"
      @name-updated="onGroupNameUpdated"
      @photo-updated="onGroupPhotoUpdated"
      @member-added="onGroupMemberAdded"
      @left-group="onGroupLeft"
    />

    <ForwardModal
      v-if="showForwardModal"
      :message="messageToForward"
      :chats="chats"
      :current-chat-id="selectedChatId"
      @close="closeForwardModal"
      @forwarded="loadChats"
    />

    <PhotoModal
      v-if="showPhotoModal"
      :is-sending="isSendingPhoto"
      @close="closePhotoModal"
      @send="sendPhotoMessage"
    />

  </div>
</template>

<script>
//Importazione componenti 
import ChatItem from '../components/ChatItem.vue'
import CreateChatModal from '../components/CreateChatModal.vue'
import CreateGroupModal from '../components/CreateGroupModal.vue'
import GroupSettingsModal from '../components/GroupSettingsModal.vue'
import ForwardModal from '../components/ForwardModal.vue'
import PhotoModal from '../components/PhotoModal.vue'

export default {
  name: "HomeView",
  components: {
    ChatItem,
    CreateChatModal,
    CreateGroupModal,
    GroupSettingsModal,
    ForwardModal,
    PhotoModal
  },
  data() {
    return {
      showCreateChat: false,
      chats: [],
      selectedChat: null,
      selectedChatId: null,
      newMessage: "",
      userId: localStorage.getItem("token"),
      showCreateGroup: false,
      showGroupSettings: false,

      // Reazioni
      activeEmojiPickerMessageId: null,
      emojiPickerPlacement: 'top',

      // Risposte
      replyingToMessage: null,

      // Inoltro
      showForwardModal: false,
      messageToForward: null,

      // Invio foto
      showPhotoModal: false,
      isSendingPhoto: false
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
    // Metodo per tornare alla lista delle chat (versione mobile)
    goBackToList() {
      this.selectedChat = null;
      this.selectedChatId = null;
      this.replyingToMessage = null;
      this.activeEmojiPickerMessageId = null;
    },

    async openChat(chatId) {
      this.selectedChatId = chatId;

      // Recupera i messaggi della chat
      try {
        const response = await this.$axios.get(`/chats/${chatId}/messages`, {
          headers: {
            Authorization: "Bearer " + localStorage.getItem("token")
          }
        });

        // Check in caso la risposta non sia valida e ordinamento
        const messages = Array.isArray(response.data) ? response.data : [];
        messages.sort((a, b) => new Date(a.instant || 0) - new Date(b.instant || 0));

        // Aggiorna chat come letta
        const currentChatInfo = this.chats.find(c => c.chatId === chatId);
        if (currentChatInfo) {
          currentChatInfo.hasUnread = false;
        }

        // Costruisce la chat selezionata
        this.selectedChat = {
          chatId,
          name: currentChatInfo?.name || "Chat",
          photo: currentChatInfo?.photo || "",
          isGroup: currentChatInfo?.isGroup,
          participants: currentChatInfo?.participants || [],
          messages
        };

        // Scrolla in fondo
        this.scrollToBottom();

      } catch (err) {
        console.log("Errore apertura chat:", err);
      }
    },

    // Formatta l'orario in formato HH:MM
    formatTime(timeString) {
      if (!timeString) return "";
      const date = new Date(timeString);
      return isNaN(date.getTime()) ? "" : date.toLocaleTimeString("it-IT", {
        hour: "2-digit",
        minute: "2-digit"
      });
    },

    // Carica tutte le chat dell'utente
    async loadChats() {
      try {
        // Chiamata all'endpoint /chats
        const response = await this.$axios.get(`/chats`, {
          headers: {
            Authorization: "Bearer " + localStorage.getItem("token")
          }
        });

        // Controllo la risposta
        const rawList = Array.isArray(response.data) ? response.data : [];

        // Mantiene solo le chat con chatId valido ed elimina eventuali duplicati
        const uniqueChats = [...new Map(rawList.filter(c => c && c.chatId).map(c => [c.chatId, c])).values()];

        // Costruisce l'array delle chat
        this.chats = uniqueChats.map(chat => {
          const isSelected = chat.chatId === this.selectedChatId;
          return {
            // Copia le proprietà della chat
            ...chat,
            hasUnread: isSelected ? false : chat.hasUnread
          };
        });

      } catch (err) {
        console.log("Errore caricamento chat:", err);
      }
    },

    // Funzione che viene richiamata quando si crea una nuova chat
    async onChatCreated(newChatId) {
      // Ricarica tutte le chat
      await this.loadChats();
      // Seleziona la nuova chat appena creata
      if (newChatId) {
        await this.openChat(newChatId);
      }
    },

    // Funzione che viene richiamata quando si invia un messaggio
    async sendMessage() {
      // Verifica che il messaggio non sia vuoto
      if (this.newMessage.trim() === "") return;

      const messageText = this.newMessage.trim();
      const payload = {
        text: messageText
      };

      // Se c'è un messaggio a cui stiamo rispondendo, lo aggiungiamo al payload
      if (this.replyingToMessage) {
        payload.replyTo = this.replyingToMessage.id;
      }

      // Invio del messaggio
      try {
        await this.$axios.post(
          `/chats/${this.selectedChatId}/messages`,
          payload,
          {
            headers: {
              Authorization: "Bearer " + localStorage.getItem("token")
            }
          }
        );

        // Reset
        this.newMessage = "";
        this.replyingToMessage = null;

        // Ricarica la chat e scrolla in fondo
        await this.openChat(this.selectedChatId);
        await this.loadChats();
        this.scrollToBottom();

      } catch (err) {
        console.log("Errore invio messaggio:", err);
      }
    },

    // Funzione che viene richiamata quando si chiude il pannello di invio foto
    closePhotoModal() {
      this.showPhotoModal = false;
      this.isSendingPhoto = false;
    },
  
    // Funzione che viene richiamata quando si invia una foto
    async sendPhotoMessage(url) {
      // Verifica che l'url non sia vuoto
      if (!url) return;

      this.isSendingPhoto = true;

      // Costruisce il payload
      const payload = {
        text: url
      };
      if (this.replyingToMessage) {
        payload.replyTo = this.replyingToMessage.id;
      }

      // Invio della foto
      try {
        await this.$axios.post(
          `/chats/${this.selectedChatId}/messages`,
          payload,
          {
            headers: {
              Authorization: "Bearer " + localStorage.getItem("token")
            }
          }
        );

        // Reset
        this.closePhotoModal();
        this.replyingToMessage = null;

        // Ricarica la chat e scrolla in fondo
        await this.openChat(this.selectedChatId);
        await this.loadChats();
        this.scrollToBottom();

      } catch (err) {
        console.error("Errore durante l'invio della foto:", err);
      } finally {
        this.isSendingPhoto = false;
      }
    },


    // Funzione che verifica se un messaggio è una foto (URL immagine o data URI)
    isMessagePhoto(text) {
      if (!text || typeof text !== "string") return false;
      const t = text.trim();
      return (
        /^(https?:\/\/.*\.(?:png|jpg|jpeg|gif|webp|svg|bmp|avif)(\?.*)?)$/i.test(t) ||
        /^data:image\//i.test(t)
      );
    },

    // Recupera le informazioni sul messaggio a cui si risponde (replyTo)
    getReplyInfo(msg) {
      if (!msg || !msg.replyTo) return null;
      if (msg._replyInfo) return msg._replyInfo;

      const originalMsg = this.selectedChat?.messages?.find(m => m.id === msg.replyTo);
      const originalText = originalMsg?.text;
      const isPhoto = this.isMessagePhoto(originalText);

      const info = {
        id: msg.replyTo,
        sender: this.getSenderName(originalMsg),
        text: originalText || "Risposta a un messaggio",
        isPhoto
      };

      msg._replyInfo = info;
      return info;
    },

    // Prende la reazione lasciata dall'utente loggato
    getUserReaction(msg) {
      return msg?.comments?.find(c => this.isMyReactionComment(c)) || null;
    },

    // Verifica se l'utente ha lasciato una reazione con emoji specifica
    hasUserReactedWith(msg, emoji) {
      return this.getUserReaction(msg)?.emoji === emoji;
    },

    // Posiziona il selettore delle emoji
    toggleEmojiPicker(messageId) {
      // Se già aperto, non fa nulla
      if (this.activeEmojiPickerMessageId === messageId) {
        this.activeEmojiPickerMessageId = null;
        this.emojiPickerPlacement = 'top';
      // Altrimenti lo apre
      } else {
        this.activeEmojiPickerMessageId = messageId;
        this.emojiPickerPlacement = 'top';

        this.$nextTick(() => {
          const container = this.$refs.messagesContainer;
          if (!container) return;
          const popup = container.querySelector('.emoji-picker-popup');
          if (!popup) return;

          // Estrae le coordinate
          const popupRect = popup.getBoundingClientRect();
          const containerRect = container.getBoundingClientRect();

          // Se non c'è abbastanza spazio sopra il messaggio (o sfora sopra), posiziona sotto al fumetto
          if (popupRect.top < containerRect.top + 8) {
            this.emojiPickerPlacement = 'bottom';
          }

          // Controllo di sicurezza orizzontale: impedisce la fuoriuscita dai bordi laterali dello schermo
          if (popupRect.right > containerRect.right - 8) {
            popup.style.right = '0px';
            popup.style.left = 'auto';
          } else if (popupRect.left < containerRect.left + 8) {
            popup.style.left = '0px';
            popup.style.right = 'auto';
          }
        });
      }
    },

    // Chiude il selettore di emoji
    closeEmojiPicker() {
      this.activeEmojiPickerMessageId = null;
      this.emojiPickerPlacement = 'top';
    },

    // Gestisce la selezione di una emoji
    async handleEmojiSelect(msg, emoji) {
      // Chiude il selettore di emoji
      this.closeEmojiPicker();

      // Controlla se l'utente ha già lasciato un emoji e in caso la rimuove
      const existing = this.getUserReaction(msg);
      if (existing) {
        await this.removeReaction(msg, existing.id);
        
        // Se l'emoji è la stessa, non fa nulla
        if (existing.emoji === emoji) return;
      }

      // Inserisce la nuova emoji
      await this.addReaction(msg, emoji);
    },

    // Aggiunge un emoji a un messaggio
    async addReaction(msg, emoji) {
      // Chiama l'API per aggiungere un emoji
      try {
        const response = await this.$axios.post(
          `/chats/${this.selectedChatId}/messages/${msg.id}/comments`,
          { emoji },
          {
            headers: {
              Authorization: "Bearer " + localStorage.getItem("token")
            }
          }
        );

        // Se il messaggio non ha ancora reazioni, inizializza l'array
        if (!msg.comments) msg.comments = [];
        // Inserisce la nuova reazione
        msg.comments.push(response.data);
      } catch (err) {
        console.error("Errore aggiunta reazione:", err);
      }
    },

    // Rimuove la reazione
    async removeReaction(msg, commentId) {
      // Chiama l'API per rimuovere la reazione
      try {
        await this.$axios.delete(
          `/chats/${this.selectedChatId}/messages/${msg.id}/comments/${commentId}`,
          {
            headers: {
              Authorization: "Bearer " + localStorage.getItem("token")
            }
          }
        );

        // Rimuove la reazione dall'array
        msg.comments = msg.comments?.filter(c => c.id !== commentId);
      } catch (err) {
        console.error("Errore rimozione reazione:", err);
      }
    },

    // Ottiene il nome del mittente di una reazione
    getReactionSenderName(comment) {
      // Se è una reazione dell'utente loggato, restituisce "Tu", altrimenti il nome del mittente
      return this.isMyReactionComment(comment) ? "Tu" : comment.senderUsername;
    },

    // Verifica se un commento è una reazione dell'utente loggato
    isMyReactionComment(comment) {
      return comment?.sender === this.userId;
    },

    // Gestisce il click su una reazione
    async handleCommentClick(msg, comment) {
      if (this.isMyReactionComment(comment)) {
        await this.removeReaction(msg, comment.id);
      }
    },

    // Inizia a rispondere a un messaggio
    startReply(msg) {
      // Imposta il messaggio a cui si sta rispondendo
      this.replyingToMessage = msg;
      // Mostra l'anteprima del messaggio
      this.$nextTick(() => {
        this.$refs.messageInput?.focus();
      });
    },

    // Annulla la risposta a un messaggio
    cancelReply() {
      this.replyingToMessage = null;
    },

    // Apre il pannello per inoltrare un messaggio
    openForwardModal(msg) {
      this.messageToForward = msg;
      this.showForwardModal = true;
    },

    // Chiude il pannello per inoltrare un messaggio
    closeForwardModal() {
      this.showForwardModal = false;
      this.messageToForward = null;
    },

    // Verifica se un messaggio è stato inviato dall'utente loggato
    isMyMessage(message) {
      return message?.sender === this.userId;
    },

    // Verifica se un messaggio è stato eliminato
    isMessageDeleted(message) {
      return message?.deleted;
    },

    // Cancella un messaggio
    async deleteMessage(msg) {
      // Chiamata API per cancellare il messaggio
      try {
        await this.$axios.delete(`/chats/${this.selectedChatId}/messages/${msg.id}`, {
          headers: {
            Authorization: "Bearer " + localStorage.getItem("token")
          }
        });

        // Aggiorna lo stato del messaggio
        msg.deleted = true;
        msg.text = "Questo messaggio è stato eliminato";
        msg.comments = [];
        msg.photo = null;
        if (this.replyingToMessage?.id === msg.id) {
          this.replyingToMessage = null;
        }

        // Aggiorna l'anteprima nella sidebar se è l'ultimo messaggio
        const chat = this.chats.find(c => c.chatId === this.selectedChatId);
        const lastMsg = this.selectedChat?.messages?.[this.selectedChat.messages.length - 1];
        if (chat && lastMsg?.id === msg.id) {
          chat.lastMessage = msg.text;
        }
      } catch (err) {
        console.error("Errore durante l'eliminazione:", err);
      }
    },

    // Restituisce il nome del mittente di un messaggio
    getSenderName(message) {
      return message?.senderUsername;
    },

    // Determina se un messaggio è stato letto da tutti i destinatari
    isMessageRead(msg) {
      return msg?.state === 2;
    },

    // Determina la classe CSS per lo stato di un messaggio
    getStatusClass(msg) {
      if (this.isMessageRead(msg)) return "status-read";
      if (msg.state === 1) return "status-received";
      return "status-sent";
    },

    // Determina il titolo per lo stato di un messaggio
    getStatusTitle(msg) {
      if (this.isMessageRead(msg)) {
        return this.selectedChat?.isGroup ? "Letto da tutti i partecipanti" : "Letto";
      }
      if (msg.state === 1) return "Consegnato / Ricevuto";
      return "Inviato";
    },

    // Funzione di polling per aggiornare automaticamente la lista dei messaggi nella chat attiva
    async pollActiveChatMessages() {
      if (!this.selectedChatId) return;

      try {
        // Chiama l'API per ottenere i messaggi della chat attiva
        const response = await this.$axios.get(`/chats/${this.selectedChatId}/messages`, {
          headers: {
            Authorization: "Bearer " + localStorage.getItem("token")
          }
        });

        // Verifica che la chat selezionata sia ancora la stessa al completamento della richiesta
        if (this.selectedChat?.chatId === this.selectedChatId) {
          // Prende i nuovi messaggi
          const newMessages = response.data || [];
          // Calcola il numero di messaggi prima dell'aggiornamento
          const oldLen = this.selectedChat.messages?.length || 0;
          // Aggiorna la lista dei messaggi
          this.selectedChat.messages = newMessages;

          // Se sono arrivati nuovi messaggi, scorre in fondo
          if (newMessages.length > oldLen) {
            this.scrollToBottom();
          }
        }
      } catch (err) {
        
      }
    },

    scrollToBottom() {
      this.$nextTick(() => {
        const container = this.$refs.messagesContainer;
        if (container) container.scrollTop = container.scrollHeight;
      });
    },

    // Aggiorna nome gruppo
    onGroupNameUpdated(newName) {
      if (this.selectedChat) this.selectedChat.name = newName;
      this.loadChats();
    },

    // Aggiorna foto gruppo
    onGroupPhotoUpdated(newPhoto) {
      if (this.selectedChat) this.selectedChat.photo = newPhoto;
      this.loadChats();
    },

    // Aggiunge membro al gruppo
    onGroupMemberAdded() {
      // Ricarica i dettagli del gruppo
      if (this.selectedChatId) this.openChat(this.selectedChatId);
      // Aggiorna la sidebar
      this.loadChats();
    },

    // Abbandona il gruppo
    onGroupLeft() {
      // Resetta la chat attiva
      this.selectedChat = null;
      this.selectedChatId = null;
      // Aggiorna la sidebar
      this.loadChats();
    }
  },

  // Viene eseguito nell momento in cui HomeView viene caricata nel browser
  mounted() {
    // Carica le conversazioni
    this.loadChats();

    // Polling periodico ogni 2.5s per aggiornare conversazioni e messaggi della chat attiva in tempo reale
    this.intervalId = setInterval(async () => {
      await this.loadChats();
      if (this.selectedChatId) {
        await this.pollActiveChatMessages();
      }
    }, 2500);
  },

  // Prima che la finestra venga chiusa elimina il timer
  beforeUnmount() {
    if (this.intervalId) {
      clearInterval(this.intervalId);
    }
  }
};
</script>

<style scoped>
/* Layout orizzontale fisso a 2 colonne su qualsiasi risoluzione */
.home-layout {
  display: flex;
  flex-direction: row;
  width: 100%;
  height: 100%;
  max-height: 100%;
  min-height: 0;
  overflow: hidden;
  background-color: #e2fff8;
  box-sizing: border-box;
  padding: 10px 14px 10px 0;
}

/* Sidebar fluida e proporzionata */
.home-sidebar {
  width: clamp(260px, 28vw, 380px);
  min-width: 240px;
  max-width: 420px;
  height: 100%;
  max-height: 100%;
  min-height: 0;
  background-color: #ffffff;
  border-right: 2px solid #ccc;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  flex-shrink: 0;

  box-shadow: 0 20px 40px rgba(0,0,0,0.18);
  border-radius: 0 22px 22px 0;
  box-sizing: border-box;
}

/* Barra superiore */
.sidebar-top-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;

  padding: 14px 16px;
  border-bottom: 2px solid #ccc;

  background-color: #ffffff;
  border-radius: 0 22px 0 0;
  flex-shrink: 0;
  gap: 8px;
  box-sizing: border-box;
}

.left-buttons {
  display: flex;
  gap: 8px;
  align-items: center;
  min-width: 0;
}

.right-button {
  display: flex;
  align-items: center;
  flex-shrink: 0;
}

/* Bottoni */
.sidebar-btn {
  padding: 10px 14px;
  background-color: #222;
  color: white;
  border: none;
  border-radius: 12px;
  font-size: 14px;
  font-weight: bold;
  cursor: pointer;
  transition: 0.25s ease;
  white-space: nowrap;
}

.sidebar-btn:hover {
  background-color: #444;
  transform: translateY(-2px);
  box-shadow: 0 8px 16px rgba(0,0,0,0.2);
}

/* LISTA CHAT */
.sidebar-content {
  flex: 1;
  overflow-y: auto;
  padding: 10px 0;
}


/* Area principale */
.home-main-area {
  flex: 1;
  min-width: 0;
  height: 100%;
  max-height: 100%;
  min-height: 0;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  background-color: #d9f7ff;
  border-radius: 22px;
  margin-left: 14px;

  box-shadow: 0 20px 40px rgba(0,0,0,0.18);
  box-sizing: border-box;
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


/* AREA CHAT */
.chat-area {
  height: 100%;
  max-height: 100%;
  min-height: 0;
  display: flex;
  flex-direction: column;
  background-color: #ffffff;
  border-radius: 22px;
  overflow: hidden;
}

/* Header: Layout Flexbox con elementi allineati perfettamente sulla stessa riga */
.chat-header {
  flex-shrink: 0;
  padding: 14px 20px;
  border-bottom: 2px solid #ccc;
  background-color: #ffffff;
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 12px;
}

.chat-header-left {
  display: flex;
  align-items: center;
  gap: 12px;
  min-width: 0;
}

.chat-header-photo {
  width: 42px;
  height: 42px;
  border-radius: 50%;
  object-fit: cover;
  border: 1px solid #ddd;
  flex-shrink: 0;
}

.chat-header-details {
  display: flex;
  flex-direction: column;
  min-width: 0;
}

.chat-header-title {
  margin: 0;
  font-size: 20px;
  font-weight: bold;
  color: #222;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.chat-header-badge {
  font-size: 11px;
  color: #0077aa;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.chat-header-actions {
  display: flex;
  align-items: center;
}

.header-settings-btn {
  background: transparent;
  border: none;
  padding: 8px;
  border-radius: 50%;
  cursor: pointer;
  color: #555;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s ease;
  font-size: 20px;
  line-height: 1;
}

.header-settings-btn:hover {
  background-color: #f0f2f5;
  color: #111;
  transform: rotate(30deg);
}

.header-settings-btn:active {
  transform: scale(0.92);
}

/* Messaggi: Container Flexbox in colonna con scroll dedicato */
.messages {
  flex: 1 1 0;
  min-height: 0; /* Fondamentale: sblocca lo scroll in un container flex */
  overflow-y: auto;
  overflow-x: hidden;
  padding: 52px 20px 20px; /* Spazio superiore per il picker del primo messaggio */
  display: flex;
  flex-direction: column;
  gap: 10px; /* Layout flessibile con spaziatura tra i messaggi come richiesto */
  background-color: #efeae2; /* Texture e sfondo chat classico WhatsApp */
  background-image: radial-gradient(#dfd9cf 1.2px, transparent 1.2px);
  background-size: 20px 20px;
}

/* Scrollbar personalizzata per l'area messaggi */
.messages::-webkit-scrollbar {
  width: 6px;
}

.messages::-webkit-scrollbar-track {
  background: transparent;
}

.messages::-webkit-scrollbar-thumb {
  background: rgba(0, 0, 0, 0.18);
  border-radius: 4px;
}

.messages::-webkit-scrollbar-thumb:hover {
  background: rgba(0, 0, 0, 0.3);
}

/* Wrapper del messaggio: allinea fumetto e reazioni */
.message-wrapper {
  display: flex;
  flex-direction: column;
  max-width: 68%;
  min-width: 100px;
  position: relative;
}

.message-wrapper.received {
  align-self: flex-start;
  align-items: flex-start;
}

.message-wrapper.sent {
  align-self: flex-end;
  align-items: flex-end;
}

/* Garanzia stacking context: il messaggio con picker attivo sta sempre in primo piano */
.message-wrapper.has-active-picker {
  z-index: 50 !important;
}

/* Nuvoletta Base del Messaggio */
.message {
  width: fit-content;
  max-width: 100%;
  min-width: 90px;
  padding: 8px 12px;
  border-radius: 12px;
  display: flex;
  flex-direction: column;
  word-wrap: break-word;
  word-break: break-word;
  box-shadow: 0 1px 2px rgba(11, 20, 26, 0.12);
  position: relative;
  transition: transform 0.15s ease, box-shadow 0.15s ease;
}

/* Messaggi Ricevuti (altri utenti): allineati a sinistra, sfondo bianco neutro */
.message-wrapper.received .message {
  align-self: flex-start;
  background-color: #ffffff;
  color: #111b21;
  border: 1px solid #e2e8f0;
  border-top-left-radius: 3px;
  border-top-right-radius: 12px;
  border-bottom-right-radius: 12px;
  border-bottom-left-radius: 12px;
}

/* Messaggi Inviati (utente corrente): allineati a destra, verde chiaro WhatsApp */
.message-wrapper.sent .message {
  align-self: flex-end;
  background-color: #dcf8c6;
  color: #111b21;
  border: 1px solid #c5e6ad;
  border-top-left-radius: 12px;
  border-top-right-radius: 3px;
  border-bottom-right-radius: 12px;
  border-bottom-left-radius: 12px;
}

/* Indicatore Inoltrato (in alto a destra del fumetto) */
.forwarded-indicator {
  display: inline-flex;
  align-items: center;
  gap: 3px;
  font-size: 11px;
  font-style: italic;
  color: #64748b;
  margin-bottom: 4px;
  align-self: flex-end;
  user-select: none;
}

/* Username mittente (visibile solo nei ricevuti) */
.message-sender {
  font-size: 12px;
  font-weight: 700;
  color: #0077aa;
  margin-bottom: 4px;
  letter-spacing: 0.2px;
}

/* Citazione Risposta incorporata nel messaggio (stile WhatsApp) */
.message-reply-preview {
  display: flex;
  background: rgba(0, 0, 0, 0.05);
  border-radius: 6px;
  margin-bottom: 6px;
  overflow: hidden;
  max-width: 100%;
}

.reply-preview-stripe {
  width: 4px;
  background-color: #0077aa;
  flex-shrink: 0;
}

.reply-preview-body {
  padding: 4px 8px;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  font-size: 12px;
  min-width: 0;
}

.reply-preview-sender {
  font-weight: 600;
  color: #0077aa;
  margin-bottom: 2px;
}

.reply-preview-text {
  color: #475569;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 280px;
}

.photo-indicator {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  color: #475569;
  font-style: italic;
}

/* Testo del messaggio */
.message-text {
  font-size: 15px;
  line-height: 1.45;
  color: #1f2937;
  white-space: pre-wrap;
}

.message-photo-wrapper {
  margin: 2px 0;
  display: flex;
  align-items: center;
}

.message-photo-img {
  max-width: 300px;
  width: 100%;
  height: auto;
  max-height: 320px;
  border-radius: 8px;
  object-fit: cover;
  display: block;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  transition: opacity 0.15s ease;
}

.message-photo-img:hover {
  opacity: 0.95;
}

/* Footer reazioni integrato dentro al fumetto del messaggio sotto al testo */
.message-reactions-footer {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
  margin-top: 6px;
  margin-bottom: 3px;
  padding-top: 5px;
  border-top: 1px solid rgba(0, 0, 0, 0.06);
}

.reaction-footer-badge {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  background-color: rgba(255, 255, 255, 0.9);
  border: 1px solid #cbd5e1;
  border-radius: 12px;
  padding: 2px 7px;
  font-size: 11px;
  cursor: pointer;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);
  transition: all 0.15s ease;
  user-select: none;
}

.reaction-footer-badge:hover {
  background-color: #ffffff;
  border-color: #94a3b8;
  transform: translateY(-1px);
}

.reaction-footer-badge.is-mine {
  background-color: #e0f2fe;
  border-color: #7dd3fc;
  color: #0369a1;
  font-weight: 600;
}

.reaction-footer-badge.is-mine:hover {
  background-color: #bae6fd;
}

.reaction-badge-emoji {
  font-size: 13px;
  line-height: 1;
}

.reaction-badge-user {
  font-size: 11px;
  color: #475569;
}

.reaction-footer-badge.is-mine .reaction-badge-user {
  color: #0284c7;
  font-weight: 600;
}

/* Dettagli footer (orario e spunte) */
.message-meta {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 4px;
  margin-top: 4px;
  user-select: none;
}

.timestamp {
  font-size: 11px;
  color: #64748b;
  font-weight: 500;
  white-space: nowrap;
}

/* Spunte di stato */
.message-status {
  display: inline-flex;
  align-items: center;
  margin-left: 2px;
  font-size: 11px;
  font-weight: 700;
  line-height: 1;
  letter-spacing: -1.5px;
}

.message-status.status-sent {
  color: #8696a0;
}

.message-status.status-received {
  color: #8696a0;
}

.message-status.status-read {
  color: #0077aa;
}

/* Barra azioni messaggio al passaggio del mouse */
.message-action-toolbar {
  position: absolute;
  top: -14px;
  right: 6px;
  display: none;
  align-items: center;
  background: #ffffff;
  border: 1px solid #cbd5e1;
  border-radius: 16px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.12);
  padding: 2px 4px;
  z-index: 20;
  gap: 2px;
}

.message-wrapper:hover .message-action-toolbar,
.message-wrapper:focus-within .message-action-toolbar {
  display: flex;
}

.has-active-picker .message-action-toolbar {
  display: none !important;
}

.msg-action-btn {
  background: none;
  border: none;
  cursor: pointer;
  padding: 3px 5px;
  border-radius: 10px;
  font-size: 13px;
  line-height: 1;
  color: #64748b;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.15s ease;
}

.msg-action-btn:hover {
  background-color: #f1f5f9;
  color: #0f172a;
}

.msg-action-btn.btn-delete-msg:hover {
  background-color: #fee2e2;
  color: #ef4444;
}

/* Messaggio eliminato (Tombstone) */
.message-wrapper.is-deleted .message {
  background-color: #f1f5f9 !important;
  color: #94a3b8 !important;
  border: 1px dashed #cbd5e1;
  box-shadow: none !important;
}

.tombstone-text {
  font-style: italic;
  color: #94a3b8;
  display: inline-flex;
  align-items: center;
  gap: 6px;
  user-select: none;
}

.tombstone-text em {
  font-style: italic;
}

/* Popup Selettore 4 Emoji: posizionato e con z-index elevato */
.emoji-picker-popup {
  position: absolute;
  top: -46px;
  z-index: 60 !important;
  display: flex;
  align-items: center;
  gap: 4px;
  background: #ffffff;
  border: 1px solid #cbd5e1;
  border-radius: 24px;
  padding: 3px 8px;
  box-shadow: 0 6px 18px rgba(0, 0, 0, 0.18);
  white-space: nowrap;
  animation: popIn 0.15s cubic-bezier(0.175, 0.885, 0.32, 1.275);
  pointer-events: auto;
}

/* Quando il messaggio è vicino al bordo superiore, posiziona il picker sotto al fumetto */
.emoji-picker-popup.placement-bottom {
  top: calc(100% + 6px);
  bottom: auto;
}

/* Messaggi Inviati (utente corrente): allineato rigorosamente a destra del fumetto per non sbordare a destra */
.message-wrapper.sent .emoji-picker-popup {
  right: 0 !important;
  left: auto !important;
}

/* Messaggi Ricevuti (altri utenti): allineato a sinistra del fumetto */
.message-wrapper.received .emoji-picker-popup {
  left: 0 !important;
  right: auto !important;
}

@keyframes popIn {
  0% { transform: scale(0.8); opacity: 0; }
  100% { transform: scale(1); opacity: 1; }
}

.emoji-option-btn {
  background: none;
  border: none;
  font-size: 20px;
  line-height: 1;
  cursor: pointer;
  padding: 4px 6px;
  border-radius: 50%;
  transition: transform 0.12s ease, background-color 0.12s ease;
}

.emoji-option-btn:hover {
  transform: scale(1.3);
  background-color: #f1f5f9;
}

.emoji-option-btn.is-selected {
  background-color: #dbeafe;
}

/* Barra Preview Risposta agganciata sopra la send-box */
.reply-preview-bar {
  display: flex;
  align-items: center;
  background-color: #f1f5f9;
  padding: 8px 14px;
  border-top: 1px solid #e2e8f0;
  border-left: 4px solid #0077aa;
  flex-shrink: 0;
  animation: slideDown 0.18s ease-out;
}

@keyframes slideDown {
  0% { transform: translateY(6px); opacity: 0; }
  100% { transform: translateY(0); opacity: 1; }
}

.reply-preview-body {
  flex: 1;
  min-width: 0;
  padding: 0 8px;
}

.reply-preview-title {
  font-size: 12px;
  color: #0077aa;
  margin-bottom: 2px;
}

.reply-preview-snippet {
  font-size: 13px;
  color: #475569;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.btn-cancel-reply {
  background: none;
  border: none;
  font-size: 18px;
  color: #64748b;
  cursor: pointer;
  padding: 4px 8px;
  border-radius: 50%;
  transition: 0.15s;
}

.btn-cancel-reply:hover {
  background-color: rgba(0, 0, 0, 0.08);
  color: #1e293b;
}

/* Nessun messaggio presente */
.no-messages {
  text-align: center;
  color: #94a3b8;
  margin: auto;
  font-size: 16px;
  font-style: italic;
  padding: 20px;
}

/* Box invio: saldamente ancorato in basso senza flex-grow */
.send-box {
  flex-shrink: 0;
  flex-grow: 0;
  display: flex;
  align-items: center;
  padding: 15px;
  border-top: 2px solid #ccc;
  background-color: #f5f5f5;
  gap: 10px;
}

.btn-attach-photo {
  background: transparent;
  border: none;
  font-size: 20px;
  line-height: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  padding: 8px;
  border-radius: 50%;
  transition: all 0.2s ease;
  flex-shrink: 0;
}

.btn-attach-photo:hover {
  background-color: #e2e8f0;
  color: #0f172a;
}

.send-box input {
  flex: 1;
  padding: 12px 14px;
  border-radius: 12px;
  border: 1px solid #aaa;
  font-size: 14px;
  outline: none;
  transition: border-color 0.2s;
}

.send-box input:focus {
  border-color: #0077aa;
}

.send-box .btn-send-message {
  margin-left: 0;
  padding: 12px 20px;
  border-radius: 12px;
  border: none;
  background-color: #222;
  color: white;
  font-weight: bold;
  cursor: pointer;
  transition: background-color 0.2s ease;
}

.send-box .btn-send-message:hover {
  background-color: #444;
}




/* RESPONSIVE: preserva sempre il layout affiancato su 2 colonne senza scroll verticale dell'intera pagina */
@media (max-width: 1024px) {
  .home-sidebar {
    width: clamp(230px, 30vw, 300px);
    min-width: 220px;
  }

  .sidebar-top-bar {
    padding: 12px 12px;
    gap: 6px;
  }

  .sidebar-btn {
    padding: 8px 10px;
    font-size: 13px;
    border-radius: 10px;
  }
}

@media (max-width: 768px) {
  .home-layout {
    padding: 6px 8px 6px 0;
  }

  .home-sidebar {
    width: clamp(200px, 34vw, 260px);
    min-width: 190px;
  }

  .home-main-area {
    margin-left: 8px;
    border-radius: 16px;
  }

  .sidebar-top-bar {
    padding: 10px 8px;
    gap: 4px;
  }

  .left-buttons {
    gap: 4px;
  }

  .sidebar-btn {
    padding: 7px 8px;
    font-size: 12px;
    border-radius: 8px;
  }

  .chat-header {
    padding: 10px 14px;
  }

  .chat-header-title {
    font-size: 17px;
  }

  .messages {
    padding: 50px 12px 16px;
    gap: 8px;
  }

  .message-wrapper {
    max-width: 88%;
  }

  .message {
    max-width: 100%;
    padding: 8px 12px;
  }

  .message-text {
    font-size: 14px;
  }

  .send-box {
    padding: 8px 12px;
    gap: 8px;
  }

  .send-box input {
    padding: 10px 12px;
    font-size: 13px;
  }

  .send-box .btn-send-message {
    padding: 10px 14px;
    font-size: 13px;
  }

  .message-sender {
    font-size: 11px;
  }

  .timestamp {
    font-size: 10px;
  }
}

/* Pulsante Indietro (visibile solo su mobile) */
.btn-back-mobile {
  display: none;
  background: none;
  border: none;
  color: #333;
  cursor: pointer;
  padding: 6px;
  border-radius: 50%;
  flex-shrink: 0;
  transition: background-color 0.15s;
  font-size: 20px;
  line-height: 1;
}

.btn-back-mobile:hover {
  background-color: #f0f0f0;
}

/* ========== MOBILE: < 600px ========== */
@media (max-width: 600px) {
  .home-layout {
    padding: 0;
    position: relative;
  }

  /* Sidebar full width */
  .home-sidebar {
    width: 100% !important;
    min-width: 0 !important;
    max-width: none !important;
    border-radius: 0;
    border-right: none;
    box-shadow: none;
  }

  /* Main area full width sovrapposta */
  .home-main-area {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    margin-left: 0;
    border-radius: 0;
    box-shadow: none;
    z-index: 10;
    display: none;
  }

  /* Quando una chat è aperta: nascondi sidebar, mostra chat */
  .home-layout.chat-open .home-sidebar {
    display: none;
  }

  .home-layout.chat-open .home-main-area {
    display: flex;
  }

  /* Mostra pulsante Indietro */
  .btn-back-mobile {
    display: flex;
    align-items: center;
    justify-content: center;
  }

  /* Sidebar top bar compatta */
  .sidebar-top-bar {
    padding: 10px 12px;
    gap: 6px;
  }

  .sidebar-btn {
    padding: 8px 10px;
    font-size: 13px;
    border-radius: 10px;
  }

  /* Chat header */
  .chat-header {
    padding: 10px 12px;
  }

  .chat-header-title {
    font-size: 16px;
  }

  .chat-header-photo {
    width: 36px;
    height: 36px;
  }

  /* Messaggi */
  .messages {
    padding: 50px 8px 16px;
    gap: 8px;
  }

  .message-action-toolbar {
    top: -16px;
    padding: 3px 6px;
    gap: 3px;
  }

  .message-wrapper {
    max-width: 92%;
  }

  .message {
    padding: 7px 10px;
  }

  .message-text {
    font-size: 14px;
  }

  .message-photo-img {
    max-width: 85vw;
    max-height: 260px;
  }

  .reply-preview-text {
    max-width: 60vw;
  }

  /* Send box */
  .send-box {
    padding: 8px 10px;
    gap: 6px;
  }

  .send-box input {
    padding: 10px 12px;
    font-size: 14px;
  }

  .send-box .btn-send-message {
    padding: 10px 14px;
    font-size: 13px;
  }

  /* Empty chat */
  .empty-chat {
    font-size: 16px;
  }

  /* Reply preview bar */
  .reply-preview-bar {
    padding: 6px 10px;
  }
}
</style>