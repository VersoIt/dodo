/// <reference types="C:/Users/Ruslan/GolandProjects/diploma/fe/node_modules/@vue/language-core/types/template-helpers.d.ts" />
/// <reference types="C:/Users/Ruslan/GolandProjects/diploma/fe/node_modules/@vue/language-core/types/props-fallback.d.ts" />
import { ref, onMounted } from 'vue';
import axios from 'axios';
import { Plus } from 'lucide-vue-next';
import { useCartStore } from '../store/cart';
const cartStore = useCartStore();
const products = ref([]);
const loading = ref(true);
const error = ref(null);
const fetchProducts = async () => {
    try {
        loading.value = true;
        const response = await axios.get('/api/catalog/products');
        products.value = response.data.products || [];
    }
    catch (err) {
        console.error('Error fetching products:', err);
        error.value = 'Failed to load products. Please try again later.';
        // Fallback data for demo if API fails
        products.value = [
            { id: '1', name: 'Margherita', description: 'Classic tomato sauce, mozzarella, and fresh basil', price: 12.99, imageUrl: 'https://images.unsplash.com/photo-1574071318508-1cdbad80ad50?q=80&w=500&auto=format&fit=crop', category: 'Classic' },
            { id: '2', name: 'Pepperoni', description: 'Tomato sauce, mozzarella, and spicy pepperoni slices', price: 14.99, imageUrl: 'https://images.unsplash.com/photo-1628840042765-356cda07504e?q=80&w=500&auto=format&fit=crop', category: 'Premium' },
            { id: '3', name: 'Hawaiian', description: 'Ham and pineapple on a bed of mozzarella and tomato sauce', price: 13.99, imageUrl: 'https://images.unsplash.com/photo-1565299624946-b28f40a0ae38?q=80&w=500&auto=format&fit=crop', category: 'Classic' },
            { id: '4', name: 'BBQ Chicken', description: 'Grilled chicken, red onions, and tangy BBQ sauce', price: 15.99, imageUrl: 'https://images.unsplash.com/photo-1513104890138-7c749659a591?q=80&w=500&auto=format&fit=crop', category: 'Premium' },
        ];
    }
    finally {
        loading.value = false;
    }
};
onMounted(() => {
    fetchProducts();
});
const __VLS_ctx = {
    ...{},
    ...{},
};
let __VLS_components;
let __VLS_intrinsics;
let __VLS_directives;
__VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
    ...{ class: "space-y-8" },
});
/** @type {__VLS_StyleScopedClasses['space-y-8']} */ ;
__VLS_asFunctionalElement1(__VLS_intrinsics.section, __VLS_intrinsics.section)({
    ...{ class: "hero min-h-[40vh] rounded-3xl overflow-hidden shadow-2xl" },
    ...{ style: {} },
});
/** @type {__VLS_StyleScopedClasses['hero']} */ ;
/** @type {__VLS_StyleScopedClasses['min-h-[40vh]']} */ ;
/** @type {__VLS_StyleScopedClasses['rounded-3xl']} */ ;
/** @type {__VLS_StyleScopedClasses['overflow-hidden']} */ ;
/** @type {__VLS_StyleScopedClasses['shadow-2xl']} */ ;
__VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
    ...{ class: "hero-overlay bg-black bg-opacity-60" },
});
/** @type {__VLS_StyleScopedClasses['hero-overlay']} */ ;
/** @type {__VLS_StyleScopedClasses['bg-black']} */ ;
/** @type {__VLS_StyleScopedClasses['bg-opacity-60']} */ ;
__VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
    ...{ class: "hero-content text-center text-neutral-content" },
});
/** @type {__VLS_StyleScopedClasses['hero-content']} */ ;
/** @type {__VLS_StyleScopedClasses['text-center']} */ ;
/** @type {__VLS_StyleScopedClasses['text-neutral-content']} */ ;
__VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
    ...{ class: "max-w-md" },
});
/** @type {__VLS_StyleScopedClasses['max-w-md']} */ ;
__VLS_asFunctionalElement1(__VLS_intrinsics.h1, __VLS_intrinsics.h1)({
    ...{ class: "mb-5 text-5xl font-extrabold uppercase tracking-tighter" },
});
/** @type {__VLS_StyleScopedClasses['mb-5']} */ ;
/** @type {__VLS_StyleScopedClasses['text-5xl']} */ ;
/** @type {__VLS_StyleScopedClasses['font-extrabold']} */ ;
/** @type {__VLS_StyleScopedClasses['uppercase']} */ ;
/** @type {__VLS_StyleScopedClasses['tracking-tighter']} */ ;
__VLS_asFunctionalElement1(__VLS_intrinsics.p, __VLS_intrinsics.p)({
    ...{ class: "mb-5 text-lg" },
});
/** @type {__VLS_StyleScopedClasses['mb-5']} */ ;
/** @type {__VLS_StyleScopedClasses['text-lg']} */ ;
__VLS_asFunctionalElement1(__VLS_intrinsics.button, __VLS_intrinsics.button)({
    ...{ class: "btn btn-primary btn-lg" },
});
/** @type {__VLS_StyleScopedClasses['btn']} */ ;
/** @type {__VLS_StyleScopedClasses['btn-primary']} */ ;
/** @type {__VLS_StyleScopedClasses['btn-lg']} */ ;
__VLS_asFunctionalElement1(__VLS_intrinsics.section, __VLS_intrinsics.section)({});
__VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
    ...{ class: "flex justify-between items-center mb-8" },
});
/** @type {__VLS_StyleScopedClasses['flex']} */ ;
/** @type {__VLS_StyleScopedClasses['justify-between']} */ ;
/** @type {__VLS_StyleScopedClasses['items-center']} */ ;
/** @type {__VLS_StyleScopedClasses['mb-8']} */ ;
__VLS_asFunctionalElement1(__VLS_intrinsics.h2, __VLS_intrinsics.h2)({
    ...{ class: "text-3xl font-bold" },
});
/** @type {__VLS_StyleScopedClasses['text-3xl']} */ ;
/** @type {__VLS_StyleScopedClasses['font-bold']} */ ;
__VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
    ...{ class: "tabs tabs-boxed" },
});
/** @type {__VLS_StyleScopedClasses['tabs']} */ ;
/** @type {__VLS_StyleScopedClasses['tabs-boxed']} */ ;
__VLS_asFunctionalElement1(__VLS_intrinsics.a, __VLS_intrinsics.a)({
    ...{ class: "tab tab-active" },
});
/** @type {__VLS_StyleScopedClasses['tab']} */ ;
/** @type {__VLS_StyleScopedClasses['tab-active']} */ ;
__VLS_asFunctionalElement1(__VLS_intrinsics.a, __VLS_intrinsics.a)({
    ...{ class: "tab" },
});
/** @type {__VLS_StyleScopedClasses['tab']} */ ;
__VLS_asFunctionalElement1(__VLS_intrinsics.a, __VLS_intrinsics.a)({
    ...{ class: "tab" },
});
/** @type {__VLS_StyleScopedClasses['tab']} */ ;
__VLS_asFunctionalElement1(__VLS_intrinsics.a, __VLS_intrinsics.a)({
    ...{ class: "tab" },
});
/** @type {__VLS_StyleScopedClasses['tab']} */ ;
if (__VLS_ctx.loading) {
    __VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
        ...{ class: "flex justify-center items-center h-64" },
    });
    /** @type {__VLS_StyleScopedClasses['flex']} */ ;
    /** @type {__VLS_StyleScopedClasses['justify-center']} */ ;
    /** @type {__VLS_StyleScopedClasses['items-center']} */ ;
    /** @type {__VLS_StyleScopedClasses['h-64']} */ ;
    __VLS_asFunctionalElement1(__VLS_intrinsics.span, __VLS_intrinsics.span)({
        ...{ class: "loading loading-spinner loading-lg text-primary" },
    });
    /** @type {__VLS_StyleScopedClasses['loading']} */ ;
    /** @type {__VLS_StyleScopedClasses['loading-spinner']} */ ;
    /** @type {__VLS_StyleScopedClasses['loading-lg']} */ ;
    /** @type {__VLS_StyleScopedClasses['text-primary']} */ ;
}
else if (__VLS_ctx.error) {
    __VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
        ...{ class: "alert alert-error" },
    });
    /** @type {__VLS_StyleScopedClasses['alert']} */ ;
    /** @type {__VLS_StyleScopedClasses['alert-error']} */ ;
    __VLS_asFunctionalElement1(__VLS_intrinsics.span, __VLS_intrinsics.span)({});
    (__VLS_ctx.error);
}
else {
    __VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
        ...{ class: "grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6" },
    });
    /** @type {__VLS_StyleScopedClasses['grid']} */ ;
    /** @type {__VLS_StyleScopedClasses['grid-cols-1']} */ ;
    /** @type {__VLS_StyleScopedClasses['md:grid-cols-2']} */ ;
    /** @type {__VLS_StyleScopedClasses['lg:grid-cols-4']} */ ;
    /** @type {__VLS_StyleScopedClasses['gap-6']} */ ;
    for (const [product] of __VLS_vFor((__VLS_ctx.products))) {
        __VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
            key: (product.id),
            ...{ class: "card bg-base-100 shadow-xl hover:shadow-2xl transition-all duration-300 group" },
        });
        /** @type {__VLS_StyleScopedClasses['card']} */ ;
        /** @type {__VLS_StyleScopedClasses['bg-base-100']} */ ;
        /** @type {__VLS_StyleScopedClasses['shadow-xl']} */ ;
        /** @type {__VLS_StyleScopedClasses['hover:shadow-2xl']} */ ;
        /** @type {__VLS_StyleScopedClasses['transition-all']} */ ;
        /** @type {__VLS_StyleScopedClasses['duration-300']} */ ;
        /** @type {__VLS_StyleScopedClasses['group']} */ ;
        __VLS_asFunctionalElement1(__VLS_intrinsics.figure, __VLS_intrinsics.figure)({
            ...{ class: "relative h-48 overflow-hidden" },
        });
        /** @type {__VLS_StyleScopedClasses['relative']} */ ;
        /** @type {__VLS_StyleScopedClasses['h-48']} */ ;
        /** @type {__VLS_StyleScopedClasses['overflow-hidden']} */ ;
        __VLS_asFunctionalElement1(__VLS_intrinsics.img)({
            src: (product.imageUrl),
            alt: (product.name),
            ...{ class: "w-full h-full object-cover group-hover:scale-110 transition-transform duration-500" },
        });
        /** @type {__VLS_StyleScopedClasses['w-full']} */ ;
        /** @type {__VLS_StyleScopedClasses['h-full']} */ ;
        /** @type {__VLS_StyleScopedClasses['object-cover']} */ ;
        /** @type {__VLS_StyleScopedClasses['group-hover:scale-110']} */ ;
        /** @type {__VLS_StyleScopedClasses['transition-transform']} */ ;
        /** @type {__VLS_StyleScopedClasses['duration-500']} */ ;
        __VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
            ...{ class: "absolute top-2 right-2" },
        });
        /** @type {__VLS_StyleScopedClasses['absolute']} */ ;
        /** @type {__VLS_StyleScopedClasses['top-2']} */ ;
        /** @type {__VLS_StyleScopedClasses['right-2']} */ ;
        __VLS_asFunctionalElement1(__VLS_intrinsics.span, __VLS_intrinsics.span)({
            ...{ class: "badge badge-secondary font-semibold" },
        });
        /** @type {__VLS_StyleScopedClasses['badge']} */ ;
        /** @type {__VLS_StyleScopedClasses['badge-secondary']} */ ;
        /** @type {__VLS_StyleScopedClasses['font-semibold']} */ ;
        (product.price);
        __VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
            ...{ class: "card-body p-6" },
        });
        /** @type {__VLS_StyleScopedClasses['card-body']} */ ;
        /** @type {__VLS_StyleScopedClasses['p-6']} */ ;
        __VLS_asFunctionalElement1(__VLS_intrinsics.h3, __VLS_intrinsics.h3)({
            ...{ class: "card-title text-xl font-bold" },
        });
        /** @type {__VLS_StyleScopedClasses['card-title']} */ ;
        /** @type {__VLS_StyleScopedClasses['text-xl']} */ ;
        /** @type {__VLS_StyleScopedClasses['font-bold']} */ ;
        (product.name);
        __VLS_asFunctionalElement1(__VLS_intrinsics.p, __VLS_intrinsics.p)({
            ...{ class: "text-sm text-base-content/70 line-clamp-2 mb-4" },
        });
        /** @type {__VLS_StyleScopedClasses['text-sm']} */ ;
        /** @type {__VLS_StyleScopedClasses['text-base-content/70']} */ ;
        /** @type {__VLS_StyleScopedClasses['line-clamp-2']} */ ;
        /** @type {__VLS_StyleScopedClasses['mb-4']} */ ;
        (product.description);
        __VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
            ...{ class: "card-actions justify-end mt-auto" },
        });
        /** @type {__VLS_StyleScopedClasses['card-actions']} */ ;
        /** @type {__VLS_StyleScopedClasses['justify-end']} */ ;
        /** @type {__VLS_StyleScopedClasses['mt-auto']} */ ;
        __VLS_asFunctionalElement1(__VLS_intrinsics.button, __VLS_intrinsics.button)({
            ...{ onClick: (...[$event]) => {
                    if (!!(__VLS_ctx.loading))
                        return;
                    if (!!(__VLS_ctx.error))
                        return;
                    __VLS_ctx.cartStore.addToCart(product);
                    // @ts-ignore
                    [loading, error, error, products, cartStore,];
                } },
            ...{ class: "btn btn-primary btn-sm gap-2" },
        });
        /** @type {__VLS_StyleScopedClasses['btn']} */ ;
        /** @type {__VLS_StyleScopedClasses['btn-primary']} */ ;
        /** @type {__VLS_StyleScopedClasses['btn-sm']} */ ;
        /** @type {__VLS_StyleScopedClasses['gap-2']} */ ;
        let __VLS_0;
        /** @ts-ignore @type {typeof __VLS_components.Plus} */
        Plus;
        // @ts-ignore
        const __VLS_1 = __VLS_asFunctionalComponent1(__VLS_0, new __VLS_0({
            ...{ class: "w-4 h-4" },
        }));
        const __VLS_2 = __VLS_1({
            ...{ class: "w-4 h-4" },
        }, ...__VLS_functionalComponentArgsRest(__VLS_1));
        /** @type {__VLS_StyleScopedClasses['w-4']} */ ;
        /** @type {__VLS_StyleScopedClasses['h-4']} */ ;
        // @ts-ignore
        [];
    }
}
// @ts-ignore
[];
const __VLS_export = (await import('vue')).defineComponent({});
export default {};
//# sourceMappingURL=HomeView.vue.js.map