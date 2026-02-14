/// <reference types="C:/Users/Ruslan/GolandProjects/diploma/fe/node_modules/@vue/language-core/types/template-helpers.d.ts" />
/// <reference types="C:/Users/Ruslan/GolandProjects/diploma/fe/node_modules/@vue/language-core/types/props-fallback.d.ts" />
import { useCartStore } from '../store/cart';
import { useRouter } from 'vue-router';
import { Trash2, Minus, Plus, ShoppingBag } from 'lucide-vue-next';
import axios from 'axios';
import { ref } from 'vue';
const cartStore = useCartStore();
const router = useRouter();
const isPlacingOrder = ref(false);
const handleCheckout = async () => {
    if (cartStore.items.length === 0)
        return;
    try {
        isPlacingOrder.value = true;
        // In a real app, we'd check if user is logged in
        const orderData = {
            items: cartStore.items.map(item => ({
                productId: item.id,
                quantity: item.quantity
            })),
            address: {
                street: 'Sample Street 123',
                city: 'Pizza City'
            }
        };
        const response = await axios.post('/api/orders', orderData);
        const orderId = response.data.id;
        cartStore.clearCart();
        router.push(`/order/${orderId}`);
    }
    catch (error) {
        console.error('Checkout failed:', error);
        alert('Failed to place order. Please try again.');
    }
    finally {
        isPlacingOrder.value = false;
    }
};
const __VLS_ctx = {
    ...{},
    ...{},
};
let __VLS_components;
let __VLS_intrinsics;
let __VLS_directives;
__VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
    ...{ class: "max-w-4xl mx-auto" },
});
/** @type {__VLS_StyleScopedClasses['max-w-4xl']} */ ;
/** @type {__VLS_StyleScopedClasses['mx-auto']} */ ;
__VLS_asFunctionalElement1(__VLS_intrinsics.h1, __VLS_intrinsics.h1)({
    ...{ class: "text-3xl font-bold mb-8" },
});
/** @type {__VLS_StyleScopedClasses['text-3xl']} */ ;
/** @type {__VLS_StyleScopedClasses['font-bold']} */ ;
/** @type {__VLS_StyleScopedClasses['mb-8']} */ ;
if (__VLS_ctx.cartStore.items.length === 0) {
    __VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
        ...{ class: "card bg-base-100 shadow-xl p-12 text-center" },
    });
    /** @type {__VLS_StyleScopedClasses['card']} */ ;
    /** @type {__VLS_StyleScopedClasses['bg-base-100']} */ ;
    /** @type {__VLS_StyleScopedClasses['shadow-xl']} */ ;
    /** @type {__VLS_StyleScopedClasses['p-12']} */ ;
    /** @type {__VLS_StyleScopedClasses['text-center']} */ ;
    __VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
        ...{ class: "flex flex-col items-center gap-4" },
    });
    /** @type {__VLS_StyleScopedClasses['flex']} */ ;
    /** @type {__VLS_StyleScopedClasses['flex-col']} */ ;
    /** @type {__VLS_StyleScopedClasses['items-center']} */ ;
    /** @type {__VLS_StyleScopedClasses['gap-4']} */ ;
    let __VLS_0;
    /** @ts-ignore @type {typeof __VLS_components.ShoppingBag} */
    ShoppingBag;
    // @ts-ignore
    const __VLS_1 = __VLS_asFunctionalComponent1(__VLS_0, new __VLS_0({
        ...{ class: "w-16 h-16 text-base-content/20" },
    }));
    const __VLS_2 = __VLS_1({
        ...{ class: "w-16 h-16 text-base-content/20" },
    }, ...__VLS_functionalComponentArgsRest(__VLS_1));
    /** @type {__VLS_StyleScopedClasses['w-16']} */ ;
    /** @type {__VLS_StyleScopedClasses['h-16']} */ ;
    /** @type {__VLS_StyleScopedClasses['text-base-content/20']} */ ;
    __VLS_asFunctionalElement1(__VLS_intrinsics.h2, __VLS_intrinsics.h2)({
        ...{ class: "text-2xl font-semibold" },
    });
    /** @type {__VLS_StyleScopedClasses['text-2xl']} */ ;
    /** @type {__VLS_StyleScopedClasses['font-semibold']} */ ;
    __VLS_asFunctionalElement1(__VLS_intrinsics.p, __VLS_intrinsics.p)({
        ...{ class: "text-base-content/60" },
    });
    /** @type {__VLS_StyleScopedClasses['text-base-content/60']} */ ;
    let __VLS_5;
    /** @ts-ignore @type {typeof __VLS_components.routerLink | typeof __VLS_components.RouterLink | typeof __VLS_components.routerLink | typeof __VLS_components.RouterLink} */
    routerLink;
    // @ts-ignore
    const __VLS_6 = __VLS_asFunctionalComponent1(__VLS_5, new __VLS_5({
        to: "/",
        ...{ class: "btn btn-primary mt-4" },
    }));
    const __VLS_7 = __VLS_6({
        to: "/",
        ...{ class: "btn btn-primary mt-4" },
    }, ...__VLS_functionalComponentArgsRest(__VLS_6));
    /** @type {__VLS_StyleScopedClasses['btn']} */ ;
    /** @type {__VLS_StyleScopedClasses['btn-primary']} */ ;
    /** @type {__VLS_StyleScopedClasses['mt-4']} */ ;
    const { default: __VLS_10 } = __VLS_8.slots;
    // @ts-ignore
    [cartStore,];
    var __VLS_8;
}
else {
    __VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
        ...{ class: "grid grid-cols-1 lg:grid-cols-3 gap-8" },
    });
    /** @type {__VLS_StyleScopedClasses['grid']} */ ;
    /** @type {__VLS_StyleScopedClasses['grid-cols-1']} */ ;
    /** @type {__VLS_StyleScopedClasses['lg:grid-cols-3']} */ ;
    /** @type {__VLS_StyleScopedClasses['gap-8']} */ ;
    __VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
        ...{ class: "lg:col-span-2 space-y-4" },
    });
    /** @type {__VLS_StyleScopedClasses['lg:col-span-2']} */ ;
    /** @type {__VLS_StyleScopedClasses['space-y-4']} */ ;
    for (const [item] of __VLS_vFor((__VLS_ctx.cartStore.items))) {
        __VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
            key: (item.id),
            ...{ class: "card card-side bg-base-100 shadow-sm border border-base-200 p-4" },
        });
        /** @type {__VLS_StyleScopedClasses['card']} */ ;
        /** @type {__VLS_StyleScopedClasses['card-side']} */ ;
        /** @type {__VLS_StyleScopedClasses['bg-base-100']} */ ;
        /** @type {__VLS_StyleScopedClasses['shadow-sm']} */ ;
        /** @type {__VLS_StyleScopedClasses['border']} */ ;
        /** @type {__VLS_StyleScopedClasses['border-base-200']} */ ;
        /** @type {__VLS_StyleScopedClasses['p-4']} */ ;
        __VLS_asFunctionalElement1(__VLS_intrinsics.figure, __VLS_intrinsics.figure)({
            ...{ class: "w-24 h-24 rounded-xl overflow-hidden flex-shrink-0" },
        });
        /** @type {__VLS_StyleScopedClasses['w-24']} */ ;
        /** @type {__VLS_StyleScopedClasses['h-24']} */ ;
        /** @type {__VLS_StyleScopedClasses['rounded-xl']} */ ;
        /** @type {__VLS_StyleScopedClasses['overflow-hidden']} */ ;
        /** @type {__VLS_StyleScopedClasses['flex-shrink-0']} */ ;
        __VLS_asFunctionalElement1(__VLS_intrinsics.img)({
            src: (item.imageUrl),
            alt: (item.name),
            ...{ class: "w-full h-full object-cover" },
        });
        /** @type {__VLS_StyleScopedClasses['w-full']} */ ;
        /** @type {__VLS_StyleScopedClasses['h-full']} */ ;
        /** @type {__VLS_StyleScopedClasses['object-cover']} */ ;
        __VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
            ...{ class: "card-body py-0 px-4 justify-between" },
        });
        /** @type {__VLS_StyleScopedClasses['card-body']} */ ;
        /** @type {__VLS_StyleScopedClasses['py-0']} */ ;
        /** @type {__VLS_StyleScopedClasses['px-4']} */ ;
        /** @type {__VLS_StyleScopedClasses['justify-between']} */ ;
        __VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
            ...{ class: "flex justify-between items-start" },
        });
        /** @type {__VLS_StyleScopedClasses['flex']} */ ;
        /** @type {__VLS_StyleScopedClasses['justify-between']} */ ;
        /** @type {__VLS_StyleScopedClasses['items-start']} */ ;
        __VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({});
        __VLS_asFunctionalElement1(__VLS_intrinsics.h3, __VLS_intrinsics.h3)({
            ...{ class: "font-bold text-lg" },
        });
        /** @type {__VLS_StyleScopedClasses['font-bold']} */ ;
        /** @type {__VLS_StyleScopedClasses['text-lg']} */ ;
        (item.name);
        __VLS_asFunctionalElement1(__VLS_intrinsics.p, __VLS_intrinsics.p)({
            ...{ class: "text-primary font-semibold" },
        });
        /** @type {__VLS_StyleScopedClasses['text-primary']} */ ;
        /** @type {__VLS_StyleScopedClasses['font-semibold']} */ ;
        (item.price);
        __VLS_asFunctionalElement1(__VLS_intrinsics.button, __VLS_intrinsics.button)({
            ...{ onClick: (...[$event]) => {
                    if (!!(__VLS_ctx.cartStore.items.length === 0))
                        return;
                    __VLS_ctx.cartStore.removeFromCart(item.id);
                    // @ts-ignore
                    [cartStore, cartStore,];
                } },
            ...{ class: "btn btn-ghost btn-sm btn-circle text-error" },
        });
        /** @type {__VLS_StyleScopedClasses['btn']} */ ;
        /** @type {__VLS_StyleScopedClasses['btn-ghost']} */ ;
        /** @type {__VLS_StyleScopedClasses['btn-sm']} */ ;
        /** @type {__VLS_StyleScopedClasses['btn-circle']} */ ;
        /** @type {__VLS_StyleScopedClasses['text-error']} */ ;
        let __VLS_11;
        /** @ts-ignore @type {typeof __VLS_components.Trash2} */
        Trash2;
        // @ts-ignore
        const __VLS_12 = __VLS_asFunctionalComponent1(__VLS_11, new __VLS_11({
            ...{ class: "w-4 h-4" },
        }));
        const __VLS_13 = __VLS_12({
            ...{ class: "w-4 h-4" },
        }, ...__VLS_functionalComponentArgsRest(__VLS_12));
        /** @type {__VLS_StyleScopedClasses['w-4']} */ ;
        /** @type {__VLS_StyleScopedClasses['h-4']} */ ;
        __VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
            ...{ class: "flex justify-between items-center" },
        });
        /** @type {__VLS_StyleScopedClasses['flex']} */ ;
        /** @type {__VLS_StyleScopedClasses['justify-between']} */ ;
        /** @type {__VLS_StyleScopedClasses['items-center']} */ ;
        __VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
            ...{ class: "join border border-base-300" },
        });
        /** @type {__VLS_StyleScopedClasses['join']} */ ;
        /** @type {__VLS_StyleScopedClasses['border']} */ ;
        /** @type {__VLS_StyleScopedClasses['border-base-300']} */ ;
        __VLS_asFunctionalElement1(__VLS_intrinsics.button, __VLS_intrinsics.button)({
            ...{ onClick: (...[$event]) => {
                    if (!!(__VLS_ctx.cartStore.items.length === 0))
                        return;
                    __VLS_ctx.cartStore.updateQuantity(item.id, item.quantity - 1);
                    // @ts-ignore
                    [cartStore,];
                } },
            ...{ class: "btn btn-ghost btn-xs join-item" },
        });
        /** @type {__VLS_StyleScopedClasses['btn']} */ ;
        /** @type {__VLS_StyleScopedClasses['btn-ghost']} */ ;
        /** @type {__VLS_StyleScopedClasses['btn-xs']} */ ;
        /** @type {__VLS_StyleScopedClasses['join-item']} */ ;
        let __VLS_16;
        /** @ts-ignore @type {typeof __VLS_components.Minus} */
        Minus;
        // @ts-ignore
        const __VLS_17 = __VLS_asFunctionalComponent1(__VLS_16, new __VLS_16({
            ...{ class: "w-3 h-3" },
        }));
        const __VLS_18 = __VLS_17({
            ...{ class: "w-3 h-3" },
        }, ...__VLS_functionalComponentArgsRest(__VLS_17));
        /** @type {__VLS_StyleScopedClasses['w-3']} */ ;
        /** @type {__VLS_StyleScopedClasses['h-3']} */ ;
        __VLS_asFunctionalElement1(__VLS_intrinsics.span, __VLS_intrinsics.span)({
            ...{ class: "btn btn-ghost btn-xs join-item pointer-events-none" },
        });
        /** @type {__VLS_StyleScopedClasses['btn']} */ ;
        /** @type {__VLS_StyleScopedClasses['btn-ghost']} */ ;
        /** @type {__VLS_StyleScopedClasses['btn-xs']} */ ;
        /** @type {__VLS_StyleScopedClasses['join-item']} */ ;
        /** @type {__VLS_StyleScopedClasses['pointer-events-none']} */ ;
        (item.quantity);
        __VLS_asFunctionalElement1(__VLS_intrinsics.button, __VLS_intrinsics.button)({
            ...{ onClick: (...[$event]) => {
                    if (!!(__VLS_ctx.cartStore.items.length === 0))
                        return;
                    __VLS_ctx.cartStore.updateQuantity(item.id, item.quantity + 1);
                    // @ts-ignore
                    [cartStore,];
                } },
            ...{ class: "btn btn-ghost btn-xs join-item" },
        });
        /** @type {__VLS_StyleScopedClasses['btn']} */ ;
        /** @type {__VLS_StyleScopedClasses['btn-ghost']} */ ;
        /** @type {__VLS_StyleScopedClasses['btn-xs']} */ ;
        /** @type {__VLS_StyleScopedClasses['join-item']} */ ;
        let __VLS_21;
        /** @ts-ignore @type {typeof __VLS_components.Plus} */
        Plus;
        // @ts-ignore
        const __VLS_22 = __VLS_asFunctionalComponent1(__VLS_21, new __VLS_21({
            ...{ class: "w-3 h-3" },
        }));
        const __VLS_23 = __VLS_22({
            ...{ class: "w-3 h-3" },
        }, ...__VLS_functionalComponentArgsRest(__VLS_22));
        /** @type {__VLS_StyleScopedClasses['w-3']} */ ;
        /** @type {__VLS_StyleScopedClasses['h-3']} */ ;
        __VLS_asFunctionalElement1(__VLS_intrinsics.p, __VLS_intrinsics.p)({
            ...{ class: "font-bold" },
        });
        /** @type {__VLS_StyleScopedClasses['font-bold']} */ ;
        ((item.price * item.quantity).toFixed(2));
        // @ts-ignore
        [];
    }
    __VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
        ...{ class: "lg:col-span-1" },
    });
    /** @type {__VLS_StyleScopedClasses['lg:col-span-1']} */ ;
    __VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
        ...{ class: "card bg-base-100 shadow-xl border border-base-200" },
    });
    /** @type {__VLS_StyleScopedClasses['card']} */ ;
    /** @type {__VLS_StyleScopedClasses['bg-base-100']} */ ;
    /** @type {__VLS_StyleScopedClasses['shadow-xl']} */ ;
    /** @type {__VLS_StyleScopedClasses['border']} */ ;
    /** @type {__VLS_StyleScopedClasses['border-base-200']} */ ;
    __VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
        ...{ class: "card-body" },
    });
    /** @type {__VLS_StyleScopedClasses['card-body']} */ ;
    __VLS_asFunctionalElement1(__VLS_intrinsics.h2, __VLS_intrinsics.h2)({
        ...{ class: "card-title mb-4" },
    });
    /** @type {__VLS_StyleScopedClasses['card-title']} */ ;
    /** @type {__VLS_StyleScopedClasses['mb-4']} */ ;
    __VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
        ...{ class: "space-y-2" },
    });
    /** @type {__VLS_StyleScopedClasses['space-y-2']} */ ;
    __VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
        ...{ class: "flex justify-between" },
    });
    /** @type {__VLS_StyleScopedClasses['flex']} */ ;
    /** @type {__VLS_StyleScopedClasses['justify-between']} */ ;
    __VLS_asFunctionalElement1(__VLS_intrinsics.span, __VLS_intrinsics.span)({});
    __VLS_asFunctionalElement1(__VLS_intrinsics.span, __VLS_intrinsics.span)({});
    (__VLS_ctx.cartStore.totalPrice.toFixed(2));
    __VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
        ...{ class: "flex justify-between" },
    });
    /** @type {__VLS_StyleScopedClasses['flex']} */ ;
    /** @type {__VLS_StyleScopedClasses['justify-between']} */ ;
    __VLS_asFunctionalElement1(__VLS_intrinsics.span, __VLS_intrinsics.span)({});
    __VLS_asFunctionalElement1(__VLS_intrinsics.span, __VLS_intrinsics.span)({
        ...{ class: "text-success" },
    });
    /** @type {__VLS_StyleScopedClasses['text-success']} */ ;
    __VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
        ...{ class: "divider" },
    });
    /** @type {__VLS_StyleScopedClasses['divider']} */ ;
    __VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
        ...{ class: "flex justify-between font-bold text-xl" },
    });
    /** @type {__VLS_StyleScopedClasses['flex']} */ ;
    /** @type {__VLS_StyleScopedClasses['justify-between']} */ ;
    /** @type {__VLS_StyleScopedClasses['font-bold']} */ ;
    /** @type {__VLS_StyleScopedClasses['text-xl']} */ ;
    __VLS_asFunctionalElement1(__VLS_intrinsics.span, __VLS_intrinsics.span)({});
    __VLS_asFunctionalElement1(__VLS_intrinsics.span, __VLS_intrinsics.span)({});
    (__VLS_ctx.cartStore.totalPrice.toFixed(2));
    __VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
        ...{ class: "card-actions mt-6" },
    });
    /** @type {__VLS_StyleScopedClasses['card-actions']} */ ;
    /** @type {__VLS_StyleScopedClasses['mt-6']} */ ;
    __VLS_asFunctionalElement1(__VLS_intrinsics.button, __VLS_intrinsics.button)({
        ...{ onClick: (__VLS_ctx.handleCheckout) },
        ...{ class: "btn btn-primary btn-block" },
        disabled: (__VLS_ctx.isPlacingOrder),
    });
    /** @type {__VLS_StyleScopedClasses['btn']} */ ;
    /** @type {__VLS_StyleScopedClasses['btn-primary']} */ ;
    /** @type {__VLS_StyleScopedClasses['btn-block']} */ ;
    if (__VLS_ctx.isPlacingOrder) {
        __VLS_asFunctionalElement1(__VLS_intrinsics.span, __VLS_intrinsics.span)({
            ...{ class: "loading loading-spinner" },
        });
        /** @type {__VLS_StyleScopedClasses['loading']} */ ;
        /** @type {__VLS_StyleScopedClasses['loading-spinner']} */ ;
    }
}
// @ts-ignore
[cartStore, cartStore, handleCheckout, isPlacingOrder, isPlacingOrder,];
const __VLS_export = (await import('vue')).defineComponent({});
export default {};
//# sourceMappingURL=CartView.vue.js.map