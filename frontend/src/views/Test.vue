<template>
  <n-layout has-sider>
    <n-layout-sider bordered>
      <n-menu
          :options="menuOptions"
          v-model:value="activeTab"
          @update:value="handleMenuSelect"
      />
    </n-layout-sider>
    <n-layout-content>
      <n-tabs
          type="card"
          v-model:value="activeTab"
          @update:value="handleTabChange"
          closable
          @close="handleTabClose"

      >
        <n-tab-pane v-for="item in tabPanes" :key="item.name" :name="item.name" :label="item.label">
          <!-- Your tab content here -->
          <div style="padding: 16px;">Content of {{ item.label }}</div>
        </n-tab-pane>
      </n-tabs>
    </n-layout-content>
  </n-layout>
</template>

<script>
import { ref } from 'vue';
import { NLayout, NLayoutSider, NLayoutContent, NMenu, NTabs, NTabPane } from 'naive-ui';
import { useRouter } from 'vue-router';

export default {
  components: {
    NLayout,
    NLayoutSider,
    NLayoutContent,
    NMenu,
    NTabs,
    NTabPane
  },
  setup() {
    const router = useRouter();
    const activeTab = ref('home');
    const menuOptions = [
      {
        label: 'Home',
        key: 'home',
        path: '/home'
      },
      {
        label: 'Profile',
        key: 'profile',
        path: "/about",

      },
      {
        label: 'Settings',
        key: 'settings',
        path: "/videoList",

      },
      {
        label: 'Settings1',
        key: 'settings1',
        path: "/setting",

      },
      {
        label: 'Settings2',
        key: 'settings2',
        path: "/test",

      }
    ];
    const tabPanes = ref([
      {
        label: 'Home',
        name: 'home'
      },
      {
        label: 'Profile',
        name: 'profile'
      },
      {
        label: 'Settings',
        name: 'settings'
      }
    ]);

    const handleMenuSelect = (key) => {
      if (!tabPanes.value.find(tabPane => tabPane.name === key)) {
        tabPanes.value.push({
          label: key.charAt(0).toUpperCase() + key.slice(1),
          name: key
        });
      }
      activeTab.value = key;
    };

    const handleTabChange = (key) => {
      activeTab.value = key;
      // router.push(key);
      console.log(key)
    };

    const handleTabClose = (key) => {
      const index = tabPanes.value.findIndex(tabPane => tabPane.name === key);
      if (index !== -1) {
        tabPanes.value.splice(index, 1);
        if (activeTab.value === key) {
          if (tabPanes.value.length) {
            activeTab.value = tabPanes.value[Math.max(0, index - 1)].name;
          } else {
            activeTab.value = '';
          }
        }
      }
    };

    return {
      activeTab,
      menuOptions,
      tabPanes,
      handleMenuSelect,
      handleTabChange,
      handleTabClose
    };
  }
};
</script>

<style>
/* Add your styles here */
</style>
