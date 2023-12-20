import Home from '../views/Home.vue'
// 务必要使用此种方法引用一个组件，否则wails打包后不能正常显示组件，我也不清楚怎么肥事儿
const routes = [
    {
        path: '/',
        component: Home
    },
    {
        name: 'user',
        path: '/user',
        component: () => import('../views/Config.vue'),
        meta: {
            keepAlive: true // 标记该页面需要缓存
        }
    },
    // {
    //     name: 'mditor',
    //     path: '/mditor',
    //     component: () => import('../views/Mditor.vue')
    // },
    // {
    //     name: 'wditor',
    //     path: '/wditor',
    //     component: () => import('../views/WangEditor.vue')
    // },
    {
        name: 'setting',
        path: '/setting',
        component: () => import('../views/Setting.vue'),
        meta: {
            keepAlive: true // 标记该页面需要缓存
        }
    },
    // {
    //     name: 'todo',
    //     path: '/todo',
    //     component: () => import('../views/Todo.vue')
    // },
    {
        name: 'about',
        path: '/about',
        component: () => import('../views/TemplateVideo.vue'),
        meta: {
            keepAlive: true // 标记该页面需要缓存
        }
    },
    {
        name: 'videoList',
        path: '/videoList',
        component: () => import('../views/VideoList.vue'),
        meta: {
            keepAlive: true // 标记该页面需要缓存
        }
    },
    {
        name: 'login',
        path: '/login',
        component: () => import('../views/Login.vue'),
        meta: {
            keepAlive: true // 标记该页面需要缓存
        }
    },
    {
        name: 'test',
        path: '/test',
        component: () => import('../views/Test.vue'),
        meta: {
            keepAlive: true // 标记该页面需要缓存
        }
    },
]

export default routes