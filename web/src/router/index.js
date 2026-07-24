import { createRouter, createWebHistory } from 'vue-router'
import Patient from '../views/Patient.vue'
import BidanApp from '../views/BidanApp.vue'

const routes = [
  { path: '/', name: 'patient', component: Patient },       // halaman publik pasien
  { path: '/app', name: 'app', component: BidanApp },       // Bidan App (privat)
]

export default createRouter({
  history: createWebHistory(),
  routes,
})
