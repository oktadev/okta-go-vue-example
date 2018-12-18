import Vue from 'vue';
import Vuex from 'vuex';

import APIClient from './apiClient';

Vue.use(Vuex);

const store = new Vuex.Store({
  state: {
    kudos: {},
    repos: [],
  },
  mutations: {
    resetRepos (state, repos) {
      state.repos = repos;
    },
    resetKudos(state, kudos) {
      state.kudos = kudos;
    }
  },
  getters: {
    allKudos(state) {
      return Object.values(state.kudos);
    },
    kudos(state) {
      return state.kudos;
    },
    repos(state) {
      return state.repos;
    },
    isKudo(state) {
      return (repo)=> {
        return !!state.kudos[repo.id];
      };
    }
  },
  actions: {
    getKudos ({commit}) {
      APIClient.getKudos().then((data) => {
        commit('resetKudos', data.reduce((acc, kudo) => { 
                               return {[kudo.id]: kudo, ...acc}
                             }, {}))
      })
    },
    updateKudo({ commit, state }, repo) {
      const kudos = { ...state.kudos, [repo.id]: repo };
      
      return APIClient
        .updateKudo(repo)
        .then(() => {
          commit('resetKudos', kudos)
        });
    },
    toggleKudo({ commit, state }, repo) {
      if (!state.kudos[repo.id]) {
        return APIClient
          .createKudo(repo)
          .then(kudo => commit('resetKudos', { [kudo.id]: kudo, ...state.kudos }))
      }

      const kudos = Object.entries(state.kudos).reduce((acc, [repoId, kudo]) => {
                      return (repoId == repo.id) ? acc
                                                 : { [repoId]: kudo, ...acc };
                    }, {});

      return APIClient
        .deleteKudo(repo)
        .then(() => commit('resetKudos', kudos));
    }
  }
});

export default store;