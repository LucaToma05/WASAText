<template>
  <div class="home-layout" :class="{ 'chat-open': selectedChat }">

    <!-- SIDEBAR -->
    <div class="home-sidebar">

      <!-- Barra dei pulsanti -->
      <div class="sidebar-top-bar">
        <div class="left-buttons">
          <button class="sidebar-btn" @click="showCreateChat = true"> Nuova chat </button>
          <button class="sidebar-btn" @click="openCreateGroupModal">Nuovo gruppo</button>
        </div>

        <div class="right-button">
          <button class="sidebar-btn" @click="goToProfile">Profilo</button>
        </div>
      </div>

      <!-- LISTA CHAT -->
      <div class="sidebar-content">
        <ChatItem
          v-for="chat in sortedChats"
          :key="chat.chatId"
          :chat="chat"
          :is-selected="chat.chatId === selectedChatId"
          :is-unread="isChatUnread(chat)"
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

      <!-- CHAT APERTA -->
      <div v-else class="chat-area">

        <!-- HEADER -->
        <div class="chat-header">
          <div class="chat-header-left">
            <button class="btn-back-mobile" @click="goBackToList" title="Torna alla lista chat">
              <svg viewBox="0 0 24 24" width="20" height="20" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <polyline points="15 18 9 12 15 6"></polyline>
              </svg>
            </button>
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
              @click="openGroupSettings" 
              title="Gestione gruppo"
              aria-label="Gestione gruppo"
            >
              <svg 
                class="gear-icon" 
                viewBox="0 0 24 24" 
                width="22" 
                height="22" 
                stroke="currentColor" 
                stroke-width="2" 
                fill="none" 
                stroke-linecap="round" 
                stroke-linejoin="round"
              >
                <circle cx="12" cy="12" r="3"></circle>
                <path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-2 2 2 2 0 0 1-2-2v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83 0 2 2 0 0 1 0-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1-2-2 2 2 0 0 1 2-2h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 0-2.83 2 2 0 0 1 2.83 0l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 2-2 2 2 0 0 1 2 2v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 0 2 2 0 0 1 0 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 2 2 2 2 0 0 1-2 2h-.09a1.65 1.65 0 0 0-1.51 1z"></path>
              </svg>
            </button>
          </div>
        </div>

        <!-- MESSAGGI -->
        <div class="messages" ref="messagesContainer" @click="closeEmojiPicker">
          <div 
            v-for="msg in selectedChat.messages"
            :key="msg.id"
            class="message-wrapper"
            :class="{ 
              'sent': isMyMessage(msg), 
              'received': !isMyMessage(msg),
              'is-mine': isMyMessage(msg), 
              'has-active-picker': activeEmojiPickerMessageId === msg.id 
            }"
          >
            <!-- Nuvoletta Messaggio -->
            <div 
              class="message"
              :class="{ 
                'sent': isMyMessage(msg), 
                'received': !isMyMessage(msg),
                'is-mine': isMyMessage(msg), 
                'is-deleted': isMessageDeleted(msg),
                'has-active-picker': activeEmojiPickerMessageId === msg.id 
              }"
            >
              <!-- Indicatore Inoltrato (in alto a destra del fumetto in stile WhatsApp) -->
              <div v-if="!isMessageDeleted(msg) && msg.forwarded" class="forwarded-indicator">
                <svg viewBox="0 0 24 24" width="12" height="12" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <polyline points="15 14 20 9 15 4"></polyline>
                  <path d="M4 20v-7a4 4 0 0 1 4-4h12"></path>
                </svg>
                <span>Inoltrato</span>
              </div>

              <!-- Mittente: visibile solo nei messaggi ricevuti -->
              <div v-if="!isMyMessage(msg)" class="message-sender">
                {{ getSenderName(msg) }}
              </div>

              <!-- Citazione Risposta incorporata (stile WhatsApp) -->
              <div v-if="!isMessageDeleted(msg) && getParsedMessage(msg).reply" class="message-reply-preview">
                <div class="reply-preview-stripe"></div>
                <div class="reply-preview-body">
                  <span class="reply-preview-sender">{{ getParsedMessage(msg).reply.sender }}</span>
                  <div class="reply-preview-text">
                    <span v-if="getParsedMessage(msg).reply.isPhoto" class="photo-indicator">
                      <svg viewBox="0 0 24 24" width="13" height="13" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                        <rect x="3" y="3" width="18" height="18" rx="2" ry="2"></rect>
                        <circle cx="8.5" cy="8.5" r="1.5"></circle>
                        <polyline points="21 15 16 10 5 21"></polyline>
                      </svg>
                      Foto
                    </span>
                    <span v-else>{{ getParsedMessage(msg).reply.text }}</span>
                  </div>
                </div>
              </div>

              <!-- Testo del messaggio / Foto / Tombstone -->
              <div class="message-text" :class="{ 'is-deleted-text': isMessageDeleted(msg) }">
                <span v-if="isMessageDeleted(msg)" class="tombstone-text">
                  <em>Questo messaggio è stato eliminato</em>
                </span>
                <div v-else-if="isMessagePhoto(getParsedMessage(msg).text)" class="message-photo-wrapper">
                  <img 
                    :src="getMessagePhotoUrl(getParsedMessage(msg).text)" 
                    class="message-photo-img" 
                    alt="Foto messaggio" 
                    loading="lazy"
                  />
                </div>
                <span v-else>{{ getParsedMessage(msg).text }}</span>
              </div>

              <!-- Footer del fumetto: Contenitore dedicato alle reazioni posizionato sotto al testo (nascosto se eliminato) -->
              <div 
                v-if="!isMessageDeleted(msg) && msg.comments && msg.comments.length > 0" 
                class="message-reactions-footer"
              >
                <button 
                  v-for="(comment, cIdx) in msg.comments" 
                  :key="comment.id || (comment.emoji + '-' + cIdx)"
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

              <!-- Dettagli: Orario e Spunte di stato -->
              <div class="message-meta">
                <span class="timestamp">{{ formatTime(msg.instant) }}</span>

                <!-- Spunte di stato: visibili solo per i messaggi inviati dall'utente e non eliminati -->
                <span 
                  v-if="isMyMessage(msg) && !isMessageDeleted(msg)" 
                  class="message-status" 
                  :class="getStatusClass(msg)"
                  :title="getStatusTitle(msg)"
                >
                  <!-- Doppia spunta: visibile SOLO se letto da tutti i destinatari (altro utente in 1-a-1 o tutti i membri nei gruppi) -->
                  <svg 
                    v-if="isMessageRead(msg)" 
                    class="check-icon double" 
                    viewBox="0 0 18 12" 
                    width="16" 
                    height="11" 
                    fill="none" 
                    xmlns="http://www.w3.org/2000/svg"
                  >
                    <path d="M1 6.5L4.5 10L11.5 2" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"/>
                    <path d="M6 6.5L9.5 10L16.5 2" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"/>
                  </svg>

                  <!-- Singola spunta: consegnato / ricevuto (in attesa di lettura da parte di tutti i destinatari) -->
                  <svg 
                    v-else 
                    class="check-icon single" 
                    viewBox="0 0 14 12" 
                    width="13" 
                    height="11" 
                    fill="none" 
                    xmlns="http://www.w3.org/2000/svg"
                  >
                    <path d="M1.5 6.5L5 10L12 2" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"/>
                  </svg>
                </span>
              </div>

              <!-- Barra azioni messaggio al passaggio del mouse: nascosta per i messaggi eliminati -->
              <div v-if="!isMessageDeleted(msg)" class="message-action-toolbar">
                <!-- Bottone Reazioni -->
                <button 
                  type="button" 
                  class="msg-action-btn" 
                  title="Reagisci"
                  @click.stop="toggleEmojiPicker(msg.id)"
                >
                  <svg viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                    <circle cx="12" cy="12" r="10"></circle>
                    <path d="M8 14s1.5 2 4 2 4-2 4-2"></path>
                    <line x1="9" y1="9" x2="9.01" y2="9"></line>
                    <line x1="15" y1="9" x2="15.01" y2="9"></line>
                  </svg>
                </button>

                <!-- Bottone Rispondi -->
                <button 
                  type="button" 
                  class="msg-action-btn" 
                  title="Rispondi"
                  @click.stop="startReply(msg)"
                >
                  <svg viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                    <polyline points="9 17 4 12 9 7"></polyline>
                    <path d="M20 18v-2a4 4 0 0 0-4-4H4"></path>
                  </svg>
                </button>

                <!-- Bottone Inoltra -->
                <button 
                  type="button" 
                  class="msg-action-btn" 
                  title="Inoltra"
                  @click.stop="openForwardModal(msg)"
                >
                  <svg viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                    <polyline points="15 17 20 12 15 7"></polyline>
                    <path d="M4 18v-2a4 4 0 0 1 4-4h12"></path>
                  </svg>
                </button>

                <!-- Bottone Elimina (visibile e cliccabile solo per i messaggi inviati dall'utente) -->
                <button 
                  v-if="isMyMessage(msg)" 
                  type="button" 
                  class="msg-action-btn btn-delete-msg" 
                  title="Elimina messaggio"
                  @click.stop.prevent="deleteMessage(msg)"
                >
                  <svg viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                    <polyline points="3 6 5 6 21 6"></polyline>
                    <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path>
                    <line x1="10" y1="11" x2="10" y2="17"></line>
                    <line x1="14" y1="11" x2="14" y2="17"></line>
                  </svg>
                </button>
              </div>

              <!-- Popup Selettore 4 Emoji (👍, ❤️, 😂, 😮) -->
              <div 
                v-if="activeEmojiPickerMessageId === msg.id" 
                class="emoji-picker-popup"
                :class="{ 'placement-bottom': emojiPickerPlacement === 'bottom' }"
                @click.stop
              >
                <button 
                  v-for="emoji in reactionEmojis" 
                  :key="emoji"
                  type="button"
                  class="emoji-option-btn"
                  :class="{ 'is-selected': hasUserReactedWith(msg, emoji) }"
                  @click="handleEmojiSelect(msg, emoji)"
                  :title="hasUserReactedWith(msg, emoji) ? 'Rimuovi reazione' : 'Reagisci con ' + emoji"
                >
                  {{ emoji }}
                </button>
              </div>
            </div>

          </div>

          <div v-if="selectedChat.messages.length === 0" class="no-messages">
            Nessun messaggio ancora. Inizia la conversazione!
          </div>
        </div>

        <!-- BARRA PREVIEW RISPOSTA (quando si risponde a un messaggio) -->
        <div v-if="replyingToMessage" class="reply-preview-bar">
          <div class="reply-preview-stripe"></div>
          <div class="reply-preview-body">
            <div class="reply-preview-title">
              Risposta a <strong>{{ getSenderName(replyingToMessage) }}</strong>
            </div>
            <div class="reply-preview-snippet">
              <span v-if="isMessagePhoto(getParsedMessage(replyingToMessage).text)" class="photo-indicator">
                📷 Foto
              </span>
              <span v-else>{{ getMessageSnippet(getParsedMessage(replyingToMessage).text, 80) }}</span>
            </div>
          </div>
          <button type="button" class="btn-cancel-reply" @click="cancelReply" title="Annulla risposta">
            ✕
          </button>
        </div>

        <!-- BOX INVIO -->
        <div class="send-box">
          <!-- Pulsante per inviare foto tramite modale -->
          <button 
            type="button" 
            class="btn-attach-photo" 
            @click="openPhotoModal" 
            title="Invia una foto"
            aria-label="Invia una foto"
          >
            <svg viewBox="0 0 24 24" width="22" height="22" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path d="M23 19a2 2 0 0 1-2 2H3a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h4l2-3h6l2 3h4a2 2 0 0 1 2 2z"></path>
              <circle cx="12" cy="13" r="4"></circle>
            </svg>
          </button>

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

    <!-- MODALI ESTRATTI -->
    <CreateChatModal
      v-if="showCreateChat"
      @close="closeCreateChat"
      @created="onChatCreated"
    />

    <CreateGroupModal
      v-if="showCreateGroup"
      @close="closeCreateGroup"
      @created="onGroupCreated"
    />

    <GroupSettingsModal
      v-if="showGroupSettings && selectedChat && selectedChat.isGroup"
      :chat="selectedChat"
      :current-user-id="userId"
      @close="closeGroupSettings"
      @name-updated="onGroupNameUpdated"
      @photo-updated="onGroupPhotoUpdated"
      @member-added="onGroupMemberAdded"
      @left-group="onGroupLeft"
    />

    <ForwardModal
      v-if="showForwardModal && messageToForward"
      :message="messageToForward"
      :chats="chats"
      :current-chat-id="selectedChatId"
      @close="closeForwardModal"
      @forwarded="onMessageForwarded"
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
      userId: localStorage.getItem("userId") || localStorage.getItem("token"),
      showCreateGroup: false,
      showGroupSettings: false,

      // Reazioni (4 emoji: 👍, ❤️, 😂, 😮)
      reactionEmojis: ['👍', '❤️', '😂', '😮'],
      activeEmojiPickerMessageId: null,
      emojiPickerPlacement: 'top',
      userCache: {},

      // Risposte
      replyingToMessage: null,

      // Inoltro
      showForwardModal: false,
      messageToForward: null,

      // Invio foto
      showPhotoModal: false,
      isSendingPhoto: false,

      // Messaggi eliminati (tombstones)
      deletedMessageIds: new Set()
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

    goBackToList() {
      this.selectedChat = null;
      this.selectedChatId = null;
      this.replyingToMessage = null;
      this.activeEmojiPickerMessageId = null;
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

        // Conserva i tombstones dei messaggi eliminati per questa chat
        const existingDeleted = (this.selectedChat && this.selectedChat.chatId === chatId) 
          ? (this.selectedChat.messages || []).filter(m => this.isMessageDeleted(m))
          : [];
        for (const delMsg of existingDeleted) {
          if (!messages.some(m => String(m.id) === String(delMsg.id))) {
            messages.push(delMsg);
          }
        }
        messages.sort((a, b) => new Date(a.instant || 0) - new Date(b.instant || 0));

        // Popola la cache degli utenti a partire dai messaggi caricati
        for (const msg of messages) {
          const sId = (msg.sender || msg.userId || msg.idSender || "").toString().trim();
          const sName = (msg.senderUsername || msg.senderName || "").toString().trim();
          if (sId && sName && sName !== sId) {
            this.userCache[sId] = sName;
          }
        }

        const currentChatInfo = this.chats.find(c => c.chatId === chatId);
        if (currentChatInfo) {
          currentChatInfo.hasUnread = false;
          currentChatInfo.unreadCount = 0;
        }

        this.selectedChat = {
          chatId,
          name: currentChatInfo?.name || "Chat",
          photo: currentChatInfo?.photo || "",
          isGroup: Boolean(currentChatInfo?.isGroup ?? currentChatInfo?.type),
          participants: currentChatInfo?.participants || [],
          messages
        };

        this.scrollToBottom();

      } catch (err) {
        console.log("Errore apertura chat:", err);
      }
    },

    formatTime(timeString) {
      if (!timeString) return "";
      let str = timeString.toString().trim();
      if (!str) return "";

      // Normalizzazione data: se proviene da SQLite "YYYY-MM-DD HH:mm:ss", aggiunge T e Z per UTC
      if (!str.includes("T") && str.includes(" ")) {
        str = str.replace(" ", "T") + "Z";
      } else if (str.includes("T") && !str.endsWith("Z") && !str.includes("+") && !str.includes("-", 10)) {
        str += "Z";
      }

      const date = new Date(str);
      if (isNaN(date.getTime())) return "";

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

        const rawList = Array.isArray(response.data) ? response.data : [];
        const seenIds = new Set();
        const uniqueChats = [];

        for (const chat of rawList) {
          const id = chat.chatId || chat.id;
          if (id && !seenIds.has(id)) {
            seenIds.add(id);
            const isSelected = id === this.selectedChatId;
            uniqueChats.push({
              ...chat,
              chatId: id,
              isGroup: Boolean(chat.isGroup ?? chat.type),
              displayName: chat.name,
              lastMessageUserId: chat.lastMessageUserId || "",
              hasUnread: isSelected ? false : Boolean(chat.hasUnread || (chat.unreadCount && chat.unreadCount > 0)),
              unreadCount: isSelected ? 0 : (chat.unreadCount || (chat.hasUnread ? 1 : 0))
            });
          }
        }

        this.chats = uniqueChats;

      } catch (err) {
        console.log("Errore caricamento chat:", err);
      }
    },

    async onChatCreated(newChatId) {
      await this.loadChats();
      if (newChatId) {
        await this.openChat(newChatId);
      }
    },

    async sendMessage() {
      if (this.newMessage.trim() === "") return;

      let messageText = this.newMessage.trim();

      // Se stiamo rispondendo a un messaggio, incorporiamo lo snippet del testo originale
      if (this.replyingToMessage) {
        const parsedOriginal = this.getParsedMessage(this.replyingToMessage);
        const isPhoto = this.isMessagePhoto(parsedOriginal.text);
        const snippet = isPhoto ? "" : this.getMessageSnippet(parsedOriginal.text, 80);

        messageText = JSON.stringify({
          replyTo: {
            id: this.replyingToMessage.id,
            sender: this.getSenderName(this.replyingToMessage),
            text: snippet,
            isPhoto
          },
          text: this.newMessage.trim()
        });
      }

      try {
        await this.$axios.post(
          `/chats/${this.selectedChatId}/messages`,
          { text: messageText },
          {
            headers: {
              Authorization: "Bearer " + localStorage.getItem("token")
            }
          }
        );

        this.newMessage = "";
        this.replyingToMessage = null;
        await this.openChat(this.selectedChatId);
        this.scrollToBottom();

      } catch (err) {
        console.log("Errore invio messaggio:", err);
      }
    },

    openPhotoModal() {
      this.showPhotoModal = true;
    },

    closePhotoModal() {
      this.showPhotoModal = false;
      this.isSendingPhoto = false;
    },

    async sendPhotoMessage(url) {
      if (!url) return;

      let messageText = url;

      // Se stiamo rispondendo a un messaggio, incorporiamo lo snippet del testo originale
      if (this.replyingToMessage) {
        const parsedOriginal = this.getParsedMessage(this.replyingToMessage);
        const isPhoto = this.isMessagePhoto(parsedOriginal.text);
        const snippet = isPhoto ? "" : this.getMessageSnippet(parsedOriginal.text, 80);

        messageText = JSON.stringify({
          replyTo: {
            id: this.replyingToMessage.id,
            sender: this.getSenderName(this.replyingToMessage),
            text: snippet,
            isPhoto
          },
          text: url
        });
      }

      this.isSendingPhoto = true;

      try {
        await this.$axios.post(
          `/chats/${this.selectedChatId}/messages`,
          {
            text: messageText,
            type: "photo",
            url: url
          },
          {
            headers: {
              Authorization: "Bearer " + localStorage.getItem("token")
            }
          }
        );

        this.closePhotoModal();
        this.replyingToMessage = null;
        await this.openChat(this.selectedChatId);
        await this.loadChats();
        this.scrollToBottom();

      } catch (err) {
        console.error("Errore durante l'invio della foto:", err);
      } finally {
        this.isSendingPhoto = false;
      }
    },

    isChatUnread(chat) {
      if (!chat) return false;
      // Se la chat è aperta, considerala già visualizzata
      if (chat.chatId === this.selectedChatId) return false;

      // Se l'ultimo messaggio è stato inviato dall'utente loggato, non è un messaggio ricevuto non letto
      const myId = this.userId || localStorage.getItem("userId") || localStorage.getItem("token");
      if (chat.lastMessageUserId && myId && String(chat.lastMessageUserId) === String(myId)) {
        return false;
      }

      return Boolean(chat.hasUnread || (chat.unreadCount && chat.unreadCount > 0));
    },

    isChatLastMessagePhoto(text) {
      if (!text) return false;
      const t = String(text).trim();
      if (t === "Questo messaggio è stato eliminato") return false;
      const parsed = this.getParsedMessage({ text });
      const target = parsed.text || text;
      if (String(target).trim() === "Questo messaggio è stato eliminato") return false;
      return this.isMessagePhoto(target);
    },

    formatChatLastMessage(text) {
      if (!text) return "";
      const t = String(text).trim();
      if (t === "Questo messaggio è stato eliminato") {
        return "Questo messaggio è stato eliminato";
      }
      const parsed = this.getParsedMessage({ text });
      const target = parsed.text || text;
      if (String(target).trim() === "Questo messaggio è stato eliminato") {
        return "Questo messaggio è stato eliminato";
      }
      if (this.isMessagePhoto(target)) {
        return "📷 Foto";
      }
      if (parsed.reply) {
        return target || "Risposta";
      }
      return target;
    },

    isMessagePhoto(text) {
      if (!text || typeof text !== "string") return false;
      const t = text.trim();
      if (t.startsWith("{") && t.endsWith("}")) {
        try {
          const parsed = JSON.parse(t);
          if (parsed.type === "photo" || parsed.isPhoto || (parsed.url && parsed.type === "photo")) {
            return true;
          }
        } catch (e) {}
      }
      return (
        /^(https?:\/\/.*\.(?:png|jpg|jpeg|gif|webp|svg|bmp|avif)(\?.*)?)$/i.test(t) ||
        /^data:image\//i.test(t) ||
        (/^(https?:\/\/[^\s]+)$/i.test(t) && (
          /\.(png|jpe?g|gif|webp|svg|bmp|avif)(\?|$)/i.test(t) ||
          /(?:image|photo|pic|avatar|placeholder|unsplash|imgur|cloudinary)/i.test(t)
        ))
      );
    },

    getMessagePhotoUrl(text) {
      if (!text || typeof text !== "string") return "";
      const t = text.trim();
      if (t.startsWith("{") && t.endsWith("}")) {
        try {
          const parsed = JSON.parse(t);
          if (parsed.url) return parsed.url;
          if (parsed.text) return parsed.text;
        } catch (e) {}
      }
      return t;
    },

    getParsedMessage(msg) {
      if (!msg || !msg.text) return { reply: null, text: "" };
      const raw = msg.text;
      if (typeof raw === "string" && (raw.startsWith('{"replyTo"') || raw.startsWith('{"_reply"'))) {
        try {
          const parsed = JSON.parse(raw);
          return {
            reply: parsed.replyTo || parsed._reply || null,
            text: parsed.text || ""
          };
        } catch (e) {
          // fallback su stringa pura
        }
      }
      return { reply: null, text: raw };
    },

    getMessageSnippet(text, maxLen = 60) {
      if (!text) return "";
      const cleaned = text.replace(/\s+/g, " ").trim();
      if (cleaned.length <= maxLen) return cleaned;
      return cleaned.slice(0, maxLen) + "...";
    },

    getUserDisplayName(userId) {
      if (!userId) return "Utente";
      const uStr = userId.toString().trim();
      const currentId = (this.userId || localStorage.getItem("userId") || localStorage.getItem("token") || "").toString().trim();
      const currentUsername = (localStorage.getItem("username") || "").toString().trim();

      if (uStr === currentId || (currentUsername && uStr === currentUsername)) {
        return "Tu";
      }

      // 1. Controlla nella cache locale degli utenti conosciuti
      if (this.userCache && this.userCache[uStr]) {
        return this.userCache[uStr];
      }

      // 2. Cerca tra i messaggi caricati della chat corrente
      if (this.selectedChat && this.selectedChat.messages) {
        const foundMsg = this.selectedChat.messages.find(
          m => (m.sender === uStr || m.userId === uStr || m.idSender === uStr) && (m.senderUsername || m.senderName)
        );
        if (foundMsg) {
          const name = foundMsg.senderUsername || foundMsg.senderName;
          this.userCache[uStr] = name;
          return name;
        }
      }

      // 3. Se è una chat 1-a-1, l'altro utente coincide con il nome della chat
      const isGroup = this.selectedChat?.isGroup || this.chats.find(c => c.chatId === this.selectedChatId)?.isGroup;
      if (!isGroup && this.selectedChat && this.selectedChat.name) {
        this.userCache[uStr] = this.selectedChat.name;
        return this.selectedChat.name;
      }

      // 4. Cerca tra i partecipanti della chat se disponibili
      if (this.selectedChat && this.selectedChat.participants) {
        const p = this.selectedChat.participants.find(
          x => (x.userId === uStr || x.id === uStr || x.name === uStr || x.username === uStr)
        );
        if (p && (p.name || p.username)) {
          const name = p.name || p.username;
          this.userCache[uStr] = name;
          return name;
        }
      }

      // 5. Cerca nella lista delle conversazioni
      const chatMatch = this.chats.find(c => c.chatId === uStr || c.name === uStr);
      if (chatMatch && chatMatch.name) {
        this.userCache[uStr] = chatMatch.name;
        return chatMatch.name;
      }

      return uStr;
    },

    getGroupedReactions(msg) {
      if (!msg.comments || !msg.comments.length) return [];
      const currentId = this.userId || localStorage.getItem("userId") || localStorage.getItem("token");
      const groups = new Map();

      for (const comment of msg.comments) {
        const emoji = comment.emoji;
        if (!groups.has(emoji)) {
          groups.set(emoji, {
            emoji,
            comments: [],
            names: [],
            hasMyReaction: false,
            myCommentId: null
          });
        }
        const g = groups.get(emoji);
        g.comments.push(comment);

        const senderId = comment.sender || comment.idSender || comment.userId;
        const isMine = senderId === currentId;
        if (isMine) {
          g.hasMyReaction = true;
          g.myCommentId = comment.id;
          g.names.push("Tu");
        } else {
          g.names.push(this.getUserDisplayName(senderId));
        }
      }

      return Array.from(groups.values());
    },

    getUserReaction(msg) {
      if (!msg || !msg.comments) return null;
      const currentId = this.userId || localStorage.getItem("userId") || localStorage.getItem("token");
      return msg.comments.find(
        c => c.sender === currentId || c.idSender === currentId || c.userId === currentId
      );
    },

    hasUserReactedWith(msg, emoji) {
      const userReaction = this.getUserReaction(msg);
      return userReaction && userReaction.emoji === emoji;
    },

    toggleEmojiPicker(messageId) {
      if (this.activeEmojiPickerMessageId === messageId) {
        this.activeEmojiPickerMessageId = null;
        this.emojiPickerPlacement = 'top';
      } else {
        this.activeEmojiPickerMessageId = messageId;
        this.emojiPickerPlacement = 'top';

        this.$nextTick(() => {
          const container = this.$refs.messagesContainer;
          if (!container) return;
          const popup = container.querySelector('.emoji-picker-popup');
          if (!popup) return;

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

    closeEmojiPicker() {
      this.activeEmojiPickerMessageId = null;
      this.emojiPickerPlacement = 'top';
    },

    async handleEmojiSelect(msg, emoji) {
      this.closeEmojiPicker();
      const existingReaction = this.getUserReaction(msg);

      // Cliccando sulla stessa reazione: la rimuove (uncommentMessage)
      if (existingReaction && existingReaction.emoji === emoji) {
        await this.removeReaction(msg, existingReaction.id);
        return;
      }

      // Se l'utente ha già una reazione diversa: la rimuove prima (massimo 1 reazione per utente)
      if (existingReaction) {
        await this.removeReaction(msg, existingReaction.id);
      }

      // Aggiunge la nuova reazione (commentMessage)
      await this.addReaction(msg, emoji);
    },

    async addReaction(msg, emoji) {
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

        const currentId = this.userId || localStorage.getItem("userId") || localStorage.getItem("token");
        const currentUsername = localStorage.getItem("username") || "Tu";
        const newComment = {
          id: response.data?.id || ("temp-" + Date.now()),
          emoji,
          sender: currentId,
          senderUsername: currentUsername,
          senderName: currentUsername,
          idSender: currentId,
          userId: currentId,
          messageId: msg.id
        };

        const list = Array.isArray(msg.comments) ? [...msg.comments] : [];
        list.push(newComment);
        msg.comments = list;
      } catch (err) {
        console.error("Errore aggiunta reazione:", err);
      }
    },

    async removeReaction(msg, commentId) {
      try {
        await this.$axios.delete(
          `/chats/${this.selectedChatId}/messages/${msg.id}/comments/${commentId}`,
          {
            headers: {
              Authorization: "Bearer " + localStorage.getItem("token")
            }
          }
        );

        if (msg.comments) {
          msg.comments = msg.comments.filter(c => c.id !== commentId);
        }
      } catch (err) {
        console.error("Errore rimozione reazione:", err);
      }
    },

    getReactionSenderName(comment) {
      if (!comment) return "Utente";

      // 1. Se coincide con l'utente attualmente loggato, mostra "Tu"
      if (this.isMyReactionComment(comment)) {
        return "Tu";
      }

      // 2. Verifica se sender è un oggetto User come da specifica OpenAPI (comment.sender.name o comment.sender.username)
      if (typeof comment.sender === "object" && comment.sender !== null) {
        if (comment.sender.name && comment.sender.name.trim() !== "") {
          return comment.sender.name;
        }
        if (comment.sender.username && comment.sender.username.trim() !== "") {
          return comment.sender.username;
        }
      }

      // 3. Verifica se l'oggetto commento ha campi username/name arricchiti
      if (comment.senderUsername && comment.senderUsername.trim() !== "" && comment.senderUsername !== comment.sender) {
        return comment.senderUsername;
      }
      if (comment.senderName && comment.senderName.trim() !== "" && comment.senderName !== comment.sender) {
        return comment.senderName;
      }
      if (comment.username && comment.username.trim() !== "" && comment.username !== comment.sender) {
        return comment.username;
      }
      if (comment.name && comment.name.trim() !== "" && comment.name !== comment.sender) {
        return comment.name;
      }

      // 4. Estrae l'identificatore del mittente
      const senderId = (
        (typeof comment.sender === "object" ? (comment.sender?.userId || comment.sender?.id) : comment.sender) ||
        comment.idSender ||
        comment.userId ||
        ""
      ).toString().trim();

      if (!senderId) return "Utente";

      // 5. Risoluzione tramite funzione di mappatura utenti
      return this.getUserDisplayName(senderId);
    },

    isMyReactionComment(comment) {
      if (!comment) return false;
      const senderId = (
        (typeof comment.sender === "object" ? (comment.sender?.userId || comment.sender?.id) : comment.sender) ||
        comment.idSender ||
        comment.userId ||
        ""
      ).toString().trim();

      const currentId = (this.userId || localStorage.getItem("userId") || localStorage.getItem("token") || "").toString().trim();
      const currentUsername = (localStorage.getItem("username") || "").toString().trim();

      if (currentId && senderId === currentId) {
        return true;
      }
      if (currentUsername) {
        const commentName = (
          (typeof comment.sender === "object" ? (comment.sender?.name || comment.sender?.username) : "") ||
          comment.senderUsername ||
          comment.senderName ||
          comment.username ||
          comment.name ||
          ""
        ).toString().trim();
        if (commentName && (commentName === currentUsername || senderId === currentUsername)) {
          return true;
        }
      }
      return false;
    },

    async handleCommentClick(msg, comment) {
      if (this.isMyReactionComment(comment)) {
        await this.removeReaction(msg, comment.id);
      }
    },

    startReply(msg) {
      this.replyingToMessage = msg;
      this.$nextTick(() => {
        this.$refs.messageInput?.focus();
      });
    },

    cancelReply() {
      this.replyingToMessage = null;
    },

    openForwardModal(msg) {
      this.messageToForward = msg;
      this.showForwardModal = true;
    },

    closeForwardModal() {
      this.showForwardModal = false;
      this.messageToForward = null;
    },

    onMessageForwarded(targetChat) {
      this.loadChats();
    },

    isMyMessage(message) {
      if (!message) return false;
      const currentId = (this.userId || localStorage.getItem("userId") || localStorage.getItem("token") || "").toString().trim();
      const currentUsername = (localStorage.getItem("username") || "").toString().trim();

      const sender = (message.sender || message.userId || message.idSender || "").toString().trim();
      const senderUsername = (message.senderUsername || message.senderName || "").toString().trim();

      if (currentId && (sender === currentId || senderUsername === currentId)) {
        return true;
      }
      if (currentUsername && (sender === currentUsername || senderUsername === currentUsername)) {
        return true;
      }
      return false;
    },

    isMessageDeleted(message) {
      if (!message) return false;
      if (message.isDeleted) return true;
      if (this.deletedMessageIds && this.deletedMessageIds.has(String(message.id))) return true;
      const text = (message.text || "").toString().trim();
      return text === "Questo messaggio è stato eliminato";
    },

    async deleteMessage(messageOrId) {
      console.log("[deleteMessage] Click evento eliminazione ricevuto:", messageOrId);

      const targetMsg = (typeof messageOrId === 'object' && messageOrId !== null)
        ? messageOrId
        : (this.selectedChat?.messages || []).find(m => String(m.id) === String(messageOrId));

      const messageId = targetMsg?.id || (typeof messageOrId === 'string' || typeof messageOrId === 'number' ? String(messageOrId) : null);
      const chatId = this.selectedChatId || (this.selectedChat && this.selectedChat.chatId) || targetMsg?.chatId;

      console.log("[deleteMessage] Dettagli parametri:", { messageId, chatId, targetMsg });

      if (!messageId) {
        console.error("[deleteMessage] Errore: messageId non trovato per", messageOrId);
        return;
      }
      if (!chatId) {
        console.error("[deleteMessage] Errore: chatId non trovato!");
        return;
      }

      try {
        const token = localStorage.getItem("token");
        console.log(`[deleteMessage] Invio richiesta API: DELETE /chats/${chatId}/messages/${messageId}`);

        const response = await this.$axios.delete(`/chats/${chatId}/messages/${messageId}`, {
          headers: {
            Authorization: "Bearer " + token
          }
        });

        console.log(`[deleteMessage] Chiamata API completata con successo. Status:`, response?.status);

        // Aggiornamento reattivo dell'interfaccia (Tombstone)
        this.markMessageAsDeleted(messageId);
      } catch (err) {
        console.error("[deleteMessage] Errore durante l'eliminazione:", err);
        if (err.response) {
          console.error("[deleteMessage] Errore risposta server:", {
            status: err.response.status,
            data: err.response.data
          });
        }
      }
    },

    deleteMessageHandler(messageOrId) {
      return this.deleteMessage(messageOrId);
    },

    markMessageAsDeleted(messageId) {
      const idStr = String(messageId);
      if (!this.deletedMessageIds) {
        this.deletedMessageIds = new Set();
      }
      this.deletedMessageIds.add(idStr);

      if (this.selectedChat && Array.isArray(this.selectedChat.messages)) {
        const idx = this.selectedChat.messages.findIndex(m => String(m.id) === idStr);
        if (idx !== -1) {
          const old = this.selectedChat.messages[idx];
          const updated = {
            ...old,
            text: "Questo messaggio è stato eliminato",
            isDeleted: true,
            comments: [],
            forwarded: false,
            photo: null,
            url: null
          };
          // Aggiornamento reattivo profondo con splice + riassegnazione array
          this.selectedChat.messages.splice(idx, 1, updated);
          this.selectedChat.messages = [...this.selectedChat.messages];
          console.log("[markMessageAsDeleted] Messaggio aggiornato reattivamente nell'array:", updated);
        } else {
          console.warn("[markMessageAsDeleted] Messaggio con id non trovato nell'array:", idStr);
        }
      }

      if (this.replyingToMessage && String(this.replyingToMessage.id) === idStr) {
        this.replyingToMessage = null;
      }

      // Se il messaggio eliminato è l'ultimo della chat selezionata, aggiorna anche la preview nella sidebar
      const currentChat = this.chats.find(c => c.chatId === this.selectedChatId);
      if (currentChat && this.selectedChat && Array.isArray(this.selectedChat.messages) && this.selectedChat.messages.length > 0) {
        const lastMsg = this.selectedChat.messages[this.selectedChat.messages.length - 1];
        if (String(lastMsg.id) === idStr) {
          currentChat.lastMessage = "Questo messaggio è stato eliminato";
        }
      }
    },

    getSenderName(message) {
      if (!message) return "";
      if (message.senderUsername && message.senderUsername.trim() !== "") {
        return message.senderUsername;
      }
      if (message.senderName && message.senderName.trim() !== "") {
        return message.senderName;
      }
      if (message.author && message.author.trim() !== "") {
        return message.author;
      }

      // Nelle chat 1-a-1, il mittente dei messaggi ricevuti è l'altro utente (nome della chat)
      const isGroup = this.selectedChat?.isGroup || this.chats.find(c => c.chatId === this.selectedChatId)?.isGroup;
      if (!isGroup && this.selectedChat && this.selectedChat.name) {
        return this.selectedChat.name;
      }

      return message.sender || "Utente";
    },

    /**
     * Determina se un messaggio inviato dall'utente corrente è stato letto da tutti i destinatari:
     * - Nelle chat 1-a-1: non appena l'altro utente visualizza il messaggio (state === 2).
     * - Nei gruppi: solo ed esclusivamente quando tutti i membri del gruppo (escluso me) lo hanno visualizzato.
     */
    isMessageRead(msg) {
      if (!msg) return false;

      // Nello standard WASAText / OpenAPI, state === 2 indica messaggio letto da tutti i destinatari
      if (msg.state === 2) {
        return true;
      }

      // Supporto se il backend fornisce la lista esplicita dei lettori (readBy / readers / readUsers)
      const rawReaders = msg.readBy || msg.readers || msg.readUsers;
      if (Array.isArray(rawReaders) && rawReaders.length > 0) {
        const myId = (this.userId || localStorage.getItem("userId") || "").toString().trim();
        const myUsername = (localStorage.getItem("username") || "").toString().trim();

        // Escludi il mittente (me stesso) dai lettori
        const otherReaders = rawReaders.filter(r => {
          const rId = (typeof r === "object" ? (r.userId || r.id || r.username || r.name) : r).toString().trim();
          return rId !== myId && rId !== myUsername;
        });

        const isGroup = Boolean(this.selectedChat?.isGroup);

        if (!isGroup) {
          // Chat 1-a-1: l'altro utente ha letto
          return otherReaders.length >= 1;
        } else {
          // Nei gruppi: tutti i membri del gruppo (escluso me) devono aver letto
          const participants = this.selectedChat?.participants || [];
          const otherParticipants = participants.filter(p => {
            const pId = (typeof p === "object" ? (p.userId || p.id || p.username || p.name) : p).toString().trim();
            return pId !== myId && pId !== myUsername;
          });

          if (otherParticipants.length > 0) {
            return otherReaders.length >= otherParticipants.length;
          }

          if (msg.totalRecipients && msg.readCount) {
            return msg.readCount >= msg.totalRecipients;
          }
        }
      }

      // Se state === 1 (Ricevuto/Consegnato) o stato 0, o non tutti hanno ancora letto, resta singola spunta
      return false;
    },

    getStatusClass(target) {
      // Supporta sia l'oggetto messaggio che il valore numerico di state
      if (typeof target === "object" && target !== null) {
        if (this.isMessageRead(target)) return "status-read";
        if (target.state === 1) return "status-received";
        return "status-sent";
      }
      if (target === 2) return "status-read";
      if (target === 1) return "status-received";
      return "status-sent";
    },

    getStatusTitle(target) {
      if (typeof target === "object" && target !== null) {
        if (this.isMessageRead(target)) {
          return this.selectedChat?.isGroup ? "Letto da tutti i partecipanti" : "Letto";
        }
        if (target.state === 1) return "Consegnato / Ricevuto";
        return "Inviato";
      }
      if (target === 2) return "Letto";
      if (target === 1) return "Ricevuto";
      return "Inviato";
    },

    async pollActiveChatMessages() {
      if (!this.selectedChatId) return;

      try {
        const response = await this.$axios.get(`/chats/${this.selectedChatId}/messages`, {
          headers: {
            Authorization: "Bearer " + localStorage.getItem("token")
          }
        });

        const newMessages = Array.isArray(response.data) ? response.data : [];

        // Verifica che la chat selezionata sia ancora la stessa al completamento della richiesta
        if (this.selectedChat && this.selectedChat.chatId === this.selectedChatId) {
          const oldLen = this.selectedChat.messages ? this.selectedChat.messages.length : 0;

          // Conserva i messaggi eliminati (tombstones)
          const existingDeleted = (this.selectedChat.messages || []).filter(m => this.isMessageDeleted(m));
          const merged = [...newMessages];
          for (const delMsg of existingDeleted) {
            if (!merged.some(m => String(m.id) === String(delMsg.id))) {
              merged.push(delMsg);
            }
          }
          merged.sort((a, b) => new Date(a.instant || 0) - new Date(b.instant || 0));

          // Assegna il nuovo array reattivamente
          this.selectedChat.messages = merged;

          // Se sono arrivati nuovi messaggi dall'esterno, scorre in fondo
          if (merged.length > oldLen) {
            this.scrollToBottom();
          }
        }
      } catch (err) {
        // Ignora silenziosamente errori transitori di polling
      }
    },

    scrollToBottom() {
      this.$nextTick(() => {
        const container = this.$refs.messagesContainer;
        if (container) {
          container.scrollTop = container.scrollHeight;
        }
      });
    },

    closeCreateChat() {
      this.showCreateChat = false;
    },

    openCreateGroupModal() {
      this.showCreateGroup = true;
    },

    closeCreateGroup() {
      this.showCreateGroup = false;
    },

    async onGroupCreated(newChatId) {
      await this.loadChats();
      if (newChatId) {
        await this.openChat(newChatId);
      }
    },

    openGroupSettings() {
      if (!this.selectedChat || !this.selectedChat.isGroup) return;
      this.showGroupSettings = true;
    },

    closeGroupSettings() {
      this.showGroupSettings = false;
    },

    onGroupNameUpdated(newName) {
      if (this.selectedChat) this.selectedChat.name = newName;
      const chatInList = this.chats.find(c => c.chatId === this.selectedChatId);
      if (chatInList) {
        chatInList.name = newName;
        chatInList.displayName = newName;
      }
      this.loadChats();
    },

    onGroupPhotoUpdated(newPhoto) {
      if (this.selectedChat) this.selectedChat.photo = newPhoto;
      const chatInList = this.chats.find(c => c.chatId === this.selectedChatId);
      if (chatInList) {
        chatInList.photo = newPhoto;
      }
      this.loadChats();
    },

    onGroupMemberAdded(username) {
      this.loadChats();
      if (this.selectedChatId) {
        this.openChat(this.selectedChatId);
      }
    },

    onGroupLeft(chatId) {
      this.selectedChat = null;
      this.selectedChatId = null;
      this.loadChats();
    }
  },

  mounted() {
    this.loadChats();

    // Polling periodico ogni 2.5s per aggiornare conversazioni e messaggi della chat attiva in tempo reale
    this.intervalId = setInterval(async () => {
      await this.loadChats();
      if (this.selectedChatId) {
        await this.pollActiveChatMessages();
      }
    }, 2500);
  },

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

.chat-preview {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 14px 16px;
  cursor: pointer;
  transition: 0.25s ease;
  border-bottom: 1px solid #eee;
  min-width: 0;
}

.chat-preview.selected {
  background-color: #d0f0ff;
  border-left: 5px solid #0077aa;
}

.chat-preview:hover {
  background-color: #f5f5f5;
}

.chat-photo {
  width: 50px;
  height: 50px;
  border-radius: 50%;
  flex-shrink: 0;
  object-fit: cover;
}

.chat-info {
  flex: 1;
  min-width: 0; /* Fondamentale: permette il troncamento dei testi nei figli flex */
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.chat-top-row {
  display: flex;
  justify-content: space-between;
  align-items: baseline;
  min-width: 0;
}

.chat-name {
  flex: 1;
  min-width: 0;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  font-weight: bold; /* Sempre in grassetto */
  font-size: 16px;
  color: #1e293b;
}

.chat-time {
  font-size: 12px;
  color: #64748b;
  flex-shrink: 0;
  margin-left: 10px;
  white-space: nowrap;
}

.chat-bottom-row,
.chat-last-message-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  min-width: 0;
  margin-top: 4px;
}

.chat-last-message {
  flex: 1;
  min-width: 0;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  font-size: 14px;
  font-weight: 400;
  color: #64748b;
  margin-top: 0;
}

/* Indicatore verde per messaggi non letti */
.unread-dot {
  width: 10px;
  height: 10px;
  min-width: 10px;
  min-height: 10px;
  border-radius: 50%;
  background-color: #25D366;
  flex-shrink: 0;
  margin-left: 8px;
  display: inline-block;
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

.message-wrapper.received,
.message-wrapper:not(.is-mine) {
  align-self: flex-start;
  align-items: flex-start;
}

.message-wrapper.sent,
.message-wrapper.is-mine {
  align-self: flex-end;
  align-items: flex-end;
}

/* Garanzia stacking context: il messaggio con picker attivo sta sempre in primo piano */
.message-wrapper.has-active-picker,
.message.has-active-picker {
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
.message.received,
.message:not(.is-mine) {
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
.message.sent,
.message.is-mine {
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

.forwarded-indicator svg {
  stroke: #64748b;
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
}

.message-status .check-icon {
  display: block;
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
.message:hover .message-action-toolbar,
.message-wrapper:focus-within .message-action-toolbar,
.message:focus-within .message-action-toolbar {
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

.msg-action-btn svg {
  pointer-events: none;
}

.msg-action-btn.btn-delete-msg:hover {
  background-color: #fee2e2;
  color: #ef4444;
}

/* Messaggio eliminato (Tombstone) */
.message.is-deleted {
  background-color: #f1f5f9 !important;
  color: #94a3b8 !important;
  border: 1px dashed #cbd5e1;
  box-shadow: none !important;
}

.message.is-deleted.sent,
.message.is-deleted.received {
  background-color: #f1f5f9 !important;
  color: #94a3b8 !important;
}

.message-text.is-deleted-text,
.tombstone-text {
  font-style: italic;
  color: #94a3b8;
  display: inline-flex;
  align-items: center;
  gap: 6px;
  user-select: none;
}

.message-text.is-deleted-text em,
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
.message.sent .emoji-picker-popup,
.message.is-mine .emoji-picker-popup,
.message-wrapper.sent .emoji-picker-popup,
.message-wrapper.is-mine .emoji-picker-popup {
  right: 0 !important;
  left: auto !important;
}

/* Messaggi Ricevuti (altri utenti): allineato a sinistra del fumetto */
.message.received .emoji-picker-popup,
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

/* Modal Inoltro */
.modal-forward-card {
  max-width: 440px;
}

.modal-card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 14px;
}

.btn-close-modal {
  background: none;
  border: none;
  font-size: 18px;
  color: #64748b;
  cursor: pointer;
  padding: 4px 8px;
  border-radius: 50%;
}

.btn-close-modal:hover {
  background-color: #f1f5f9;
  color: #1e293b;
}

.forward-msg-preview-card {
  background: #f8fafc;
  border: 1px solid #e2e8f0;
  border-radius: 8px;
  padding: 10px 12px;
  margin-bottom: 12px;
}

.forward-msg-preview-label {
  font-size: 11px;
  font-weight: 600;
  color: #64748b;
  text-transform: uppercase;
  margin-bottom: 4px;
}

.forward-msg-preview-content {
  font-size: 13px;
  color: #1e293b;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.forward-chats-list {
  max-height: 220px;
  overflow-y: auto;
  margin: 12px 0;
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.forward-chat-row {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 8px 10px;
  border-radius: 10px;
  cursor: pointer;
  transition: background-color 0.15s;
}

.forward-chat-row:hover {
  background-color: #f1f5f9;
}

.forward-chat-avatar {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  object-fit: cover;
}

.forward-chat-meta {
  flex: 1;
  min-width: 0;
  display: flex;
  align-items: center;
  gap: 8px;
}

.forward-chat-title {
  font-size: 14px;
  font-weight: 500;
  color: #1e293b;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.no-chats-found {
  text-align: center;
  padding: 16px;
  font-size: 13px;
  color: #94a3b8;
  font-style: italic;
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

.chat-last-photo {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  color: #475569;
  font-weight: 500;
}

.chat-last-photo .photo-icon {
  font-size: 13px;
  line-height: 1;
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
  color: #64748b;
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


.chat-last-deleted {
  font-style: italic;
  color: #94a3b8;
}

.chat-last-deleted em {
  font-style: italic;
  color: #94a3b8;
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
}

.message-sender {
  font-size: 11px;
}

.timestamp {
  font-size: 10px;
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