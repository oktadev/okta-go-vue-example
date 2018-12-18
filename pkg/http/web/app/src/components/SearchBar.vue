<template>
    <v-toolbar dark color="teal">
      <v-spacer></v-spacer>
      <v-text-field
        solo-inverted
        flat
        hide-details
        label="Search for your OOS project on Github + Press Enter"
        prepend-inner-icon="search"
        v-model="query"
        @keyup.enter="onSearchSubmition"
      ></v-text-field>
      <v-spacer></v-spacer>
      <button @click.prevent="logout">Logout</button>
    </v-toolbar>
</template>

<script>
export default {
    data() {
      return {
        query: null,
      };
    },
    props: ['defaultQuery'],
    methods: {
      onSearchSubmition() {
        this.$emit('search-submitted', this.query);
      },
      async logout () {
        await this.$auth.logout()
        this.$router.push('/')
    }
  }
}
</script>