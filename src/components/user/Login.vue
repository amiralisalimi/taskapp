<template>
    <div>
        <h2>Login</h2>
        <form @submit.prevent="submitForm">
            <div>
                <label>Username:</label>
                <input type="text" v-model="username" />
            </div>
            <div>
                <label>Password:</label>
                <input type="password" v-model="password" />
            </div>
            <div>
                <button type="submit">Login</button>
            </div>
        </form>
    </div>
</template>

<script>
export default {
    data() {
        return {
            username: '',
            password: '',
        };
    },
    methods: {
        submitForm() {
            this.$store.dispatch('login', {
                username: this.username,
                password: this.password,
            })
                .then(() => {
                    this.$router.push('/dashboard');
                })
                .catch(error => {
                    console.log(error);
                });
        },
    },
};
</script>