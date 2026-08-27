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
                // Controllo se l'username è vuoto
                if (this.username.trim() === "") {
                    this.errorMsg = "Username vuoto";
                    return;
                }

                try {
                    // Esegue una richiesta HTTP POST all'endpoint /login con il seguente body
                    const response = await this.$axios.post("/session", { 
                        name: this.username
                    });
                    
                    // Memorizza l'ID restituito dal backend
                    const userId = response.data.id; 

                    // Salvo l'ID utente
                    localStorage.setItem("userId", userId);

                    // Uso il router per andare alla home dopo il login
                    this.$router.push("/home");
                } catch (err) {
                    console.log("Errore login:", err.response?.data || err);
                    this.errorMsg = err.response?.data?.error || "Errore durante il login";
                }
            }

        }
    }
</script>

<template>
    <div class="login-wrapper">
        <div class="login-container">
            <h1 id="login-label">LOGIN</h1>
            <form @submit.prevent="login">
                <input v-model="username" type="text" placeholder="Username">
                <button type="submit">Accedi</button>
            </form>
            <p v-if="errorMsg"> {{ errorMsg }}</p>
        </div>
    </div>
</template>


<style>
/* Wrapper: centra tutto e lo sposta più in alto */
.login-wrapper {
    height: 100vh;
    display: flex;
    justify-content: center;
    align-items: flex-start;
    padding-top: 120px;
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
    animation: fadeIn 0.8s ease-out;
}

/* Titolo molto grande e elegante */
#login-label {
    font-size: 72px;
    letter-spacing: 4px;
    font-weight: bold;
    color: #222;
}

/* Spazio tra input e bottone */
.login-container form {
    display: flex;
    flex-direction: column;
    gap: 20px;   /* spazio tra input e bottone */
}

/* Input e bottone molto più grandi */
.login-container input,
.login-container button {
    width: 420px;
    height: 75px;
    font-size: 28px;
    padding: 14px;
    border-radius: 14px;
}

/* Input moderno */
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

/* Bottone premium */
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

/* Animazione card */
@keyframes fadeIn {
    from { opacity: 0; transform: translateY(40px); }
    to   { opacity: 1; transform: translateY(0); }
}

/* --- VERSIONE MOBILE --- */
@media (max-width: 600px) {

    .login-container {
        padding: 40px;
        gap: 25px;
        border-radius: 16px;
    }

    #login-label {
        font-size: 42px;
        letter-spacing: 2px;
    }

    .login-container input,
    .login-container button {
        width: 100%;
        max-width: 300px;
        height: 55px;
        font-size: 20px;
        padding: 10px;
    }
}


</style>
