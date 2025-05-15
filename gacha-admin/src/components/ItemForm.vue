<template>
  <div>
    <h2>{{ isEditMode ? 'アイテム編集' : 'アイテム追加' }}</h2>
    <form @submit.prevent="saveItem" enctype="multipart/form-data">
      <div>
        <label for="name">名前:</label>
        <input type="text" id="name" v-model="form.name" required>
      </div>
      <div>
        <label for="rarity">レアリティ:</label>
        <input type="text" id="rarity" v-model="form.rarity" required>
      </div>
      <div>
        <label for="details">詳細:</label>
        <textarea id="details" v-model="form.details"></textarea>
      </div>
      <div>
        <label for="percentage">確率 (%):</label>
        <input type="number" id="percentage" v-model.number="form.percentage" min="0" max="100" required>
      </div>
      <div>
        <label for="image">画像ファイル:</label>
        <input type="file" id="image" @change="handleFileChange">
        <img :src="previewImage" alt="プレビュー" v-if="previewImage" style="max-width: 100px; max-height: 100px;">
      </div>
      <button type="submit">{{ isEditMode ? '更新' : '追加' }}</button>
      <button type="button" @click="$emit('cancel')">キャンセル</button>
    </form>
  </div>
</template>

<script>
import axios from 'axios';

export default {
  props: {
    item: {
      type: Object,
      default: () => ({ id: null, name: '', rarity: '', details: '', percentage: 0, imageIdentifier: '' }),
    },
  },
  emits: ['item-saved', 'cancel'],
  data() {
    return {
      form: { ...this.item },
      isEditMode: this.item.id !== null,
      selectedFile: null,
      previewImage: this.item.imageIdentifier ? `/api/image?id=${this.item.imageIdentifier}` : null,
    };
  },
  watch: {
    item(newItem) {
      this.form = { ...newItem };
      this.isEditMode = newItem.id !== null;
      this.previewImage = newItem.imageIdentifier ? `/api/image?id=${newItem.imageIdentifier}` : null;
      this.selectedFile = null;
    },
  },
  methods: {
    handleFileChange(event) {
      this.selectedFile = event.target.files[0];
      if (this.selectedFile) {
        this.previewImage = URL.createObjectURL(this.selectedFile);
      } else {
        this.previewImage = this.item.imageIdentifier ? `/api/image?id=${this.item.imageIdentifier}` : null;
      }
    },
    async saveItem() {
      try {
        const formData = new FormData();
        for (const key in this.form) {
          formData.append(key, this.form[key]);
        }
        if (this.selectedFile) {
          formData.append('image', this.selectedFile, this.selectedFile.name);
        }

        let response;
        if (this.isEditMode) {
          response = await axios.put(`http://localhost:8080/api/admin/items/${this.form.id}`, formData, {
            headers: {
              'Content-Type': 'multipart/form-data',
            },
          });
        } else {
          response = await axios.post('http://localhost:8080/api/admin/items', formData, {
            headers: {
              'Content-Type': 'multipart/form-data',
            },
          });
        }
        this.$emit('item-saved', response.data);
        this.$emit('cancel');
      } catch (error) {
        console.error('アイテムの保存に失敗しました', error);
        alert('アイテムの保存に失敗しました。');
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

form {
  display: flex;
  flex-direction: column;
  max-width: 400px; /* 少し幅を広げる */
  margin-top: 1rem;
}

form > div {
  margin-bottom: 0.8rem; /* 少し間隔を広げる */
}

label {
  display: block;
  margin-bottom: 0.3rem;
  font-weight: bold;
  font-size: 0.9rem;
}

input[type="text"],
input[type="number"],
textarea {
  width: calc(100% - 0.6rem); /* 左右のpaddingを考慮 */
  padding: 0.5rem;
  border: 1px solid #ccc;
  font-size: 0.9rem;
}

textarea {
  min-height: 80px;
}

button {
  padding: 0.6rem 1rem;
  margin-top: 1rem;
  cursor: pointer;
  font-size: 0.9rem;
}

button:first-child {
  background-color: #4CAF50;
  color: white;
  border: none;
}

button:last-child {
  background-color: #f44336;
  color: white;
  border: none;
}

/* スマートフォン向けのスタイル */
@media screen and (max-width: 600px) {
  form {
    max-width: 100%; /* 幅を画面いっぱいに */
  }

  input[type="text"],
  input[type="number"],
  textarea,
  button {
    font-size: 1rem; /* フォントサイズを少し大きく */
    padding: 0.7rem;
  }

  button {
    margin-top: 0.8rem;
  }
}
</style>