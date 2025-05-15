<template>
  <div id="app">
    <h1>ガチャ管理画面</h1>
    <button @click="isAdding = true">アイテム追加</button>

    <ItemForm
      v-if="isAdding"
      @item-saved="handleItemSaved"
      @cancel="isAdding = false"
    />

    <ItemForm
      v-if="isEditing"
      :item="editingItem"
      @item-saved="handleItemSaved"
      @cancel="cancelEdit"
    />

    <ItemList
      :items="items"
      @edit-item="editItem"
      @item-deleted="handleItemDeleted"
    />
  </div>
</template>

<script>
import axios from 'axios';
import ItemList from './components/ItemList.vue';
import ItemForm from './components/ItemForm.vue';

export default {
  components: {
    ItemList,
    ItemForm,
  },
  data() {
    return {
      items: [],
      isAdding: false,
      isEditing: false, // 編集状態を管理する新しいプロパティ
      editingItem: null, // 編集対象のアイテムを格納するプロパティ
    };
  },
  async mounted() {
    await this.fetchItems();
  },
  methods: {
    async fetchItems() {
      try {
        const response = await axios.get('http://localhost:8080/api/admin/items');
        this.items = response.data;
      } catch (error) {
        console.error('アイテムの取得に失敗しました', error);
        alert('アイテムの取得に失敗しました。');
      }
    },
    editItem(item) {
      this.editingItem = { ...item };
      this.isAdding = false;
      this.isEditing = true; // 編集状態を true に設定
    },
    async handleItemSaved(savedItem) {
      if (this.editingItem) {
        // 更新の場合
        const index = this.items.findIndex(item => item.id === savedItem.id);
        if (index !== -1) {
          this.items.splice(index, 1, savedItem);
        }
        this.isEditing = false; // 編集完了後に編集状態を false に戻す
        this.editingItem = null;
      } else {
        // 追加の場合
        this.items.push(savedItem);
        this.isAdding = false;
      }
    },
    cancelEdit() {
      this.isEditing = false; // 編集キャンセル時に編集状態を false に戻す
      this.editingItem = null;
    },
    handleItemDeleted(deletedId) {
      this.items = this.items.filter(item => item.id !== deletedId);
    },
  },
};
</script>

<style>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 2rem;
  padding: 1rem;
}

h1 {
  margin-bottom: 1.5rem;
  font-size: 1.5rem;
}

button {
  padding: 0.7rem 1.2rem;
  margin: 0.3rem;
  cursor: pointer;
  font-size: 1rem;
  background-color: #007bff;
  color: white;
  border: none;
  border-radius: 5px;
}

button:hover {
  background-color: #0056b3;
}

/* スマートフォン向けのスタイル */
@media screen and (max-width: 600px) {
  h1 {
    font-size: 1.8rem;
  }

  button {
    font-size: 1.1rem;
    padding: 0.8rem 1.5rem;
  }
}
</style>