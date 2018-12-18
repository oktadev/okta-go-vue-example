<template>
  <v-app id="inspire">
    <v-content>
      <v-container fluid fill-height>
        <v-layout align-center justify-center>
          <v-flex xs12 sm8 md4>
            <v-card class="elevation-12">
              <v-toolbar dark color="teal">
                <v-toolbar-title justify-center>Login</v-toolbar-title>
              </v-toolbar>
              <v-card-text>
                <v-btn @click.prevent="login" color="primary">Sign in with Okta</v-btn>
              </v-card-text>
            </v-card>
          </v-flex>
        </v-layout>
      </v-container>
    </v-content>
  </v-app>
</template>

<script>
export default {
  data() {
    return {};
  },
  async mounted() {
    const isAuthenticated = await this.$auth.isAuthenticated();
    isAuthenticated && this.$router.push('/me');
  },
  methods: {
    login () {
      this.$auth.loginRedirect('/me')
    }
  }
}
</script>
