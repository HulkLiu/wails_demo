import { onMounted, onUnmounted, ref } from 'vue';
import { isBrowser } from '../env/is-browser';
export const useBrowserLocation = (customWindow = isBrowser ? window : null) => {
    const getWindowLocation = () => {
        const { hash, host, hostname, href, origin, pathname, port, protocol, search } = (customWindow === null || customWindow === void 0 ? void 0 : customWindow.location) || {};
        return {
            hash,
            host,
            hostname,
            href,
            origin,
            pathname,
            port,
            protocol,
            search
        };
    };
    const updateLocation = () => {
        locationState.value = getWindowLocation();
    };
    const locationState = ref(getWindowLocation());
    onMounted(() => {
        if (customWindow) {
            customWindow.addEventListener('popstate', updateLocation);
            customWindow.addEventListener('hashchange', updateLocation);
        }
    });
    onUnmounted(() => {
        if (customWindow) {
            customWindow.removeEventListener('popstate', updateLocation);
            customWindow.removeEventListener('hashchange', updateLocation);
        }
    });
    return locationState;
};
