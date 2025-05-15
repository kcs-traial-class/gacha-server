<template>
  <div>
    <h2>アイテム一覧</h2>
    <table v-if="items.length > 0">
      <thead>
        <tr>
          <th>ID</th>
          <th>名前</th>
          <th>レアリティ</th>
          <th>詳細</th>
          <th>確率 (%)</th>
          <th>画像</th>
          <th>操作</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="item in items" :key="item.id">
          <td>{{ item.id }}</td>
          <td>{{ item.name }}</td>
          <td>{{ item.rarity }}</td>
          <td>{{ item.details }}</td>
          <td>{{ item.percentage }}</td>
          <td>
            <img :src="item.image_identifier ? `/api/image?id=${item.image_identifier}` : ''" alt="アイテム画像" v-if="item.image_identifier" style="max-width: 50px; max-height: 50px;">
            <span v-else>なし</span>
          </td>
          <td>
            <button @click="$emit('edit-item', item)">編集</button>
            <button @click="deleteItem(item.id)">削除</button>
          </td>
        </tr>
      </tbody>
    </table>
    <p v-else>アイテムは登録されていません。</p>
  </div>
</template>

<script>
import axios from 'axios';

export default {
  props: ['items'],
  emits: ['edit-item', 'item-deleted'],
  methods: {
    async deleteItem(id) {
      if (confirm(`ID ${id} のアイテムを削除しますか？`)) {
        try {
          await axios.delete(`http://localhost:8080/api/admin/items/${id}`);
          this.$emit('item-deleted', id);
        } catch (error) {
          console.error('アイテムの削除に失敗しました', error);
          alert('アイテムの削除に失敗しました。');
        }
      }
    },
  },
};
</script>

<style scoped>
div {
  padding: 1rem;
}

h2 {
  margin-bottom: 1rem;
}

table {
  width: 100%;
  border-collapse: collapse;
  margin-bottom: 1rem;
  font-size: 0.9rem; /* デフォルトのフォントサイズを少し小さく */
}

th,
td {
  border: 1px solid #ccc;
  padding: 0.5rem;
  text-align: left;
}

th {
  background-color: #f0f0f0;
}

button {
  padding: 0.3rem 0.6rem;
  margin-right: 0.3rem;
  cursor: pointer;
  font-size: 0.8rem;
}

/* スマートフォン向けのスタイル */
@media screen and (max-width: 600px) {
  table {
    font-size: 0.8rem;
  }

  th,
  td {
    padding: 0.3rem;
  }

  /* 横スクロールできるようにする */
  div {
    overflow-x: auto;
  }

  table {
    width: auto;
  }

  th:nth-child(n+4), /* 詳細と確率以降のヘッダーを非表示 */
  td:nth-child(n+4) { /* 詳細と確率以降のセルを非表示 */
    display: none;
  }
}
</style>