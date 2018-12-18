<template>
  <v-card >
    <v-card-title primary-title>
      <div class="repo-card-content">
        <h3 class="headline mb-0">
          <router-link :to="{ name: 'repo-details', params: { id: repo.id }}" >{{repo.full_name}}</router-link>
        </h3>
        <div>{{repo.description}}</div>
      </div>
    </v-card-title>
    <v-card-actions>
      <v-chip>
        {{repo.language}}
      </v-chip>
      <v-spacer></v-spacer>
      <v-btn @click.prevent="toggleKudo(repo)"  flat icon color="pink">
        <v-icon v-if="isKudo(repo)">favorite</v-icon>
        <v-icon v-else>favorite_border</v-icon>
      </v-btn>
    </v-card-actions>
  </v-card>
</template>

<script>
import { mapActions } from 'vuex';

export default {
  data() {
    return {}
  },
  props: ['repo'],
  methods: {
    isKudo(repo) {
      return this.$store.getters.isKudo(repo);
    },
    ...mapActions(['toggleKudo'])
  }
}
</script>

<style>
 .repo-card-content {
   height: 90px;
   overflow: scroll;
 }
</style>

