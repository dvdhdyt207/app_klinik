import { defineStore } from 'pinia'
import { ref } from 'vue'
import { api } from '../api/client'

// Menyimpan siapa yang sedang masuk.
//
// Ini hanya untuk tampilan — penjaga sesungguhnya ada di server. Frontend tidak
// pernah memutuskan siapa yang boleh melihat data; ia cuma menghindarkan bidan
// dari layar kosong dan mengarahkannya ke halaman masuk.
export const useAuth = defineStore('auth', () => {
  const user = ref(null)
  const sudahDicek = ref(false)

  // cek menanyakan status sesi ke server. Dipanggil sekali dari penjaga rute;
  // hasilnya diingat supaya perpindahan layar berikutnya tidak menanya lagi.
  async function cek() {
    if (sudahDicek.value) return user.value
    try {
      user.value = await api.me()
    } catch {
      user.value = null
    } finally {
      sudahDicek.value = true
    }
    return user.value
  }

  async function masuk(username, password) {
    user.value = await api.login({ username, password })
    sudahDicek.value = true
    return user.value
  }

  async function keluar() {
    try {
      await api.logout()
    } finally {
      // Keadaan lokal dibersihkan apa pun jawaban server: kalau permintaan
      // keluar gagal terkirim, menahan layar tetap "masuk" justru menyesatkan.
      user.value = null
      sudahDicek.value = true
    }
  }

  // lupakan dipakai saat server menjawab 401 di tengah pemakaian, supaya
  // pemeriksaan berikutnya benar-benar bertanya lagi ke server.
  function lupakan() {
    user.value = null
    sudahDicek.value = false
  }

  return { user, sudahDicek, cek, masuk, keluar, lupakan }
})
