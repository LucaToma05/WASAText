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
                    const response = await this.$axios.post("/login", { 
                        name: this.username
                    });
                    
                    // Memorizza l'ID restituito dal backend
                    const userId = response.data.id; 

                    // Salvo l'ID utente
                    localStorage.setItem("userId", userId);

                    // Uso il router èer andare alla home dopo il login
                    this.$router.push("/");
                } catch (err) {
                    this.errorMsg = "Errore durante il login";
                }
            }

        }
    }
</script>

<template>

    <h1>Login</h1>
    <form @submit.prevent="login">
        <input v-model="username" type="text" placeholder="Username">
        <button type="submit">Login</button>
    </form>
    <p v-if="errorMsg"> {{ errorMsg }}</p>
    

</template>

<style>
</style>