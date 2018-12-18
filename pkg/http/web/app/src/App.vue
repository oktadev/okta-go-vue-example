<template>
  <v-app>
    <router-view></router-view>
    <Footer />
  </v-app>
</template>

<script>
// import Home from '@/components/Home.vue'
// import Login from '@/components/Login.vue'
import Footer from '@/components/Footer.vue'

export default {
  name: 'App',
  components: { Footer },
  data() {
    return {}
  }
}
</script>

