/// <reference types="C:/Users/Ruslan/GolandProjects/diploma/fe/node_modules/@vue/language-core/types/template-helpers.d.ts" />
/// <reference types="C:/Users/Ruslan/GolandProjects/diploma/fe/node_modules/@vue/language-core/types/props-fallback.d.ts" />
import { ref, onMounted, onUnmounted } from 'vue';
import { useRoute } from 'vue-router';
import axios from 'axios';
import { CheckCircle2, Clock, Truck, PackageCheck, CreditCard } from 'lucide-vue-next';
const route = useRoute();
const orderId = route.params.id;
const order = ref(null);
const loading = ref(true);
let pollInterval = null;
const fetchOrder = async () => {
    try {
        const response = await axios.get(`/api/orders/${orderId}`);
        order.value = response.data;
    }
    catch (error) {
        console.error('Failed to fetch order:', error);
        // Mock for demo
        if (!order.value) {
            order.value = {
                id: orderId,
                status: 0, // Received
                totalPrice: 25.98,
                items: [
                    { name: 'Pepperoni', quantity: 1, price: 14.99 },
                    { name: 'Margherita', quantity: 1, price: 12.99 }
                ]
            };
        }
    }
    finally {
        loading.value = false;
    }
};
const getStatusIndex = (status) => {
    // 0: Received, 1: Kitchen, 2: Ready, 3: InDelivery, 4: Delivered, 5: Cancelled
    return status;
};
const steps = [
    { label: 'Order Received', icon: PackageCheck },
    { label: 'In the Kitchen', icon: Clock },
    { label: 'Ready for Pickup', icon: CheckCircle2 },
    { label: 'Out for Delivery', icon: Truck },
    { label: 'Enjoy your Pizza!', icon: CheckCircle2 }
];
const handlePayment = async () => {
    try {
        await axios.post(`/api/orders/${orderId}/pay`);
        fetchOrder();
    }
    catch (error) {
        alert('Payment failed');
    }
};
onMounted(() => {
    fetchOrder();
    pollInterval = setInterval(fetchOrder, 5000);
});
onUnmounted(() => {
    if (pollInterval)
        clearInterval(pollInterval);
});
const __VLS_ctx = {
    ...{},
    ...{},
};
let __VLS_components;
let __VLS_intrinsics;
let __VLS_directives;
__VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
    ...{ class: "max-w-3xl mx-auto py-8" },
});
/** @type {__VLS_StyleScopedClasses['max-w-3xl']} */ ;
/** @type {__VLS_StyleScopedClasses['mx-auto']} */ ;
/** @type {__VLS_StyleScopedClasses['py-8']} */ ;
if (__VLS_ctx.loading) {
    __VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
        ...{ class: "flex justify-center py-20" },
    });
    /** @type {__VLS_StyleScopedClasses['flex']} */ ;
    /** @type {__VLS_StyleScopedClasses['justify-center']} */ ;
    /** @type {__VLS_StyleScopedClasses['py-20']} */ ;
    __VLS_asFunctionalElement1(__VLS_intrinsics.span, __VLS_intrinsics.span)({
        ...{ class: "loading loading-spinner loading-lg text-primary" },
    });
    /** @type {__VLS_StyleScopedClasses['loading']} */ ;
    /** @type {__VLS_StyleScopedClasses['loading-spinner']} */ ;
    /** @type {__VLS_StyleScopedClasses['loading-lg']} */ ;
    /** @type {__VLS_StyleScopedClasses['text-primary']} */ ;
}
else if (__VLS_ctx.order) {
    __VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
        ...{ class: "space-y-8" },
    });
    /** @type {__VLS_StyleScopedClasses['space-y-8']} */ ;
    __VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
        ...{ class: "flex justify-between items-center" },
    });
    /** @type {__VLS_StyleScopedClasses['flex']} */ ;
    /** @type {__VLS_StyleScopedClasses['justify-between']} */ ;
    /** @type {__VLS_StyleScopedClasses['items-center']} */ ;
    __VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({});
    __VLS_asFunctionalElement1(__VLS_intrinsics.h1, __VLS_intrinsics.h1)({
        ...{ class: "text-3xl font-bold" },
    });
    /** @type {__VLS_StyleScopedClasses['text-3xl']} */ ;
    /** @type {__VLS_StyleScopedClasses['font-bold']} */ ;
    (__VLS_ctx.orderId.slice(0, 8));
    __VLS_asFunctionalElement1(__VLS_intrinsics.p, __VLS_intrinsics.p)({
        ...{ class: "text-base-content/60" },
    });
    /** @type {__VLS_StyleScopedClasses['text-base-content/60']} */ ;
    __VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
        ...{ class: "badge badge-lg badge-primary font-bold" },
    });
    /** @type {__VLS_StyleScopedClasses['badge']} */ ;
    /** @type {__VLS_StyleScopedClasses['badge-lg']} */ ;
    /** @type {__VLS_StyleScopedClasses['badge-primary']} */ ;
    /** @type {__VLS_StyleScopedClasses['font-bold']} */ ;
    (__VLS_ctx.steps[__VLS_ctx.getStatusIndex(__VLS_ctx.order.status)]?.label || 'Processing');
    if (__VLS_ctx.order.status === 0) {
        __VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
            ...{ class: "card bg-secondary text-secondary-content shadow-xl" },
        });
        /** @type {__VLS_StyleScopedClasses['card']} */ ;
        /** @type {__VLS_StyleScopedClasses['bg-secondary']} */ ;
        /** @type {__VLS_StyleScopedClasses['text-secondary-content']} */ ;
        /** @type {__VLS_StyleScopedClasses['shadow-xl']} */ ;
        __VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
            ...{ class: "card-body flex-row items-center gap-6" },
        });
        /** @type {__VLS_StyleScopedClasses['card-body']} */ ;
        /** @type {__VLS_StyleScopedClasses['flex-row']} */ ;
        /** @type {__VLS_StyleScopedClasses['items-center']} */ ;
        /** @type {__VLS_StyleScopedClasses['gap-6']} */ ;
        let __VLS_0;
        /** @ts-ignore @type {typeof __VLS_components.CreditCard} */
        CreditCard;
        // @ts-ignore
        const __VLS_1 = __VLS_asFunctionalComponent1(__VLS_0, new __VLS_0({
            ...{ class: "w-12 h-12" },
        }));
        const __VLS_2 = __VLS_1({
            ...{ class: "w-12 h-12" },
        }, ...__VLS_functionalComponentArgsRest(__VLS_1));
        /** @type {__VLS_StyleScopedClasses['w-12']} */ ;
        /** @type {__VLS_StyleScopedClasses['h-12']} */ ;
        __VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
            ...{ class: "flex-grow" },
        });
        /** @type {__VLS_StyleScopedClasses['flex-grow']} */ ;
        __VLS_asFunctionalElement1(__VLS_intrinsics.h2, __VLS_intrinsics.h2)({
            ...{ class: "card-title text-2xl" },
        });
        /** @type {__VLS_StyleScopedClasses['card-title']} */ ;
        /** @type {__VLS_StyleScopedClasses['text-2xl']} */ ;
        __VLS_asFunctionalElement1(__VLS_intrinsics.p, __VLS_intrinsics.p)({});
        __VLS_asFunctionalElement1(__VLS_intrinsics.button, __VLS_intrinsics.button)({
            ...{ onClick: (__VLS_ctx.handlePayment) },
            ...{ class: "btn btn-neutral btn-lg" },
        });
        /** @type {__VLS_StyleScopedClasses['btn']} */ ;
        /** @type {__VLS_StyleScopedClasses['btn-neutral']} */ ;
        /** @type {__VLS_StyleScopedClasses['btn-lg']} */ ;
    }
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
        ...{ class: "card-title mb-6 text-xl" },
    });
    /** @type {__VLS_StyleScopedClasses['card-title']} */ ;
    /** @type {__VLS_StyleScopedClasses['mb-6']} */ ;
    /** @type {__VLS_StyleScopedClasses['text-xl']} */ ;
    __VLS_asFunctionalElement1(__VLS_intrinsics.ul, __VLS_intrinsics.ul)({
        ...{ class: "steps steps-vertical md:steps-horizontal w-full" },
    });
    /** @type {__VLS_StyleScopedClasses['steps']} */ ;
    /** @type {__VLS_StyleScopedClasses['steps-vertical']} */ ;
    /** @type {__VLS_StyleScopedClasses['md:steps-horizontal']} */ ;
    /** @type {__VLS_StyleScopedClasses['w-full']} */ ;
    for (const [step, index] of __VLS_vFor((__VLS_ctx.steps))) {
        __VLS_asFunctionalElement1(__VLS_intrinsics.li, __VLS_intrinsics.li)({
            key: (index),
            ...{ class: "step" },
            ...{ class: ({ 'step-primary': index <= __VLS_ctx.getStatusIndex(__VLS_ctx.order.status) }) },
        });
        /** @type {__VLS_StyleScopedClasses['step']} */ ;
        /** @type {__VLS_StyleScopedClasses['step-primary']} */ ;
        __VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
            ...{ class: "flex flex-col items-center gap-1" },
        });
        /** @type {__VLS_StyleScopedClasses['flex']} */ ;
        /** @type {__VLS_StyleScopedClasses['flex-col']} */ ;
        /** @type {__VLS_StyleScopedClasses['items-center']} */ ;
        /** @type {__VLS_StyleScopedClasses['gap-1']} */ ;
        const __VLS_5 = (step.icon);
        // @ts-ignore
        const __VLS_6 = __VLS_asFunctionalComponent1(__VLS_5, new __VLS_5({
            ...{ class: "w-6 h-6 mb-1" },
        }));
        const __VLS_7 = __VLS_6({
            ...{ class: "w-6 h-6 mb-1" },
        }, ...__VLS_functionalComponentArgsRest(__VLS_6));
        /** @type {__VLS_StyleScopedClasses['w-6']} */ ;
        /** @type {__VLS_StyleScopedClasses['h-6']} */ ;
        /** @type {__VLS_StyleScopedClasses['mb-1']} */ ;
        __VLS_asFunctionalElement1(__VLS_intrinsics.span, __VLS_intrinsics.span)({
            ...{ class: "text-xs font-semibold" },
        });
        /** @type {__VLS_StyleScopedClasses['text-xs']} */ ;
        /** @type {__VLS_StyleScopedClasses['font-semibold']} */ ;
        (step.label);
        // @ts-ignore
        [loading, order, order, order, order, orderId, steps, steps, getStatusIndex, getStatusIndex, handlePayment,];
    }
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
        ...{ class: "divide-y divide-base-200" },
    });
    /** @type {__VLS_StyleScopedClasses['divide-y']} */ ;
    /** @type {__VLS_StyleScopedClasses['divide-base-200']} */ ;
    for (const [item] of __VLS_vFor((__VLS_ctx.order.items))) {
        __VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
            key: (item.name),
            ...{ class: "py-3 flex justify-between" },
        });
        /** @type {__VLS_StyleScopedClasses['py-3']} */ ;
        /** @type {__VLS_StyleScopedClasses['flex']} */ ;
        /** @type {__VLS_StyleScopedClasses['justify-between']} */ ;
        __VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({});
        __VLS_asFunctionalElement1(__VLS_intrinsics.span, __VLS_intrinsics.span)({
            ...{ class: "font-bold" },
        });
        /** @type {__VLS_StyleScopedClasses['font-bold']} */ ;
        (item.quantity);
        (item.name);
        __VLS_asFunctionalElement1(__VLS_intrinsics.span, __VLS_intrinsics.span)({
            ...{ class: "font-semibold" },
        });
        /** @type {__VLS_StyleScopedClasses['font-semibold']} */ ;
        ((item.price * item.quantity).toFixed(2));
        // @ts-ignore
        [order,];
    }
    __VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
        ...{ class: "pt-4 mt-4 flex justify-between text-xl font-bold" },
    });
    /** @type {__VLS_StyleScopedClasses['pt-4']} */ ;
    /** @type {__VLS_StyleScopedClasses['mt-4']} */ ;
    /** @type {__VLS_StyleScopedClasses['flex']} */ ;
    /** @type {__VLS_StyleScopedClasses['justify-between']} */ ;
    /** @type {__VLS_StyleScopedClasses['text-xl']} */ ;
    /** @type {__VLS_StyleScopedClasses['font-bold']} */ ;
    __VLS_asFunctionalElement1(__VLS_intrinsics.span, __VLS_intrinsics.span)({});
    __VLS_asFunctionalElement1(__VLS_intrinsics.span, __VLS_intrinsics.span)({});
    (__VLS_ctx.order.totalPrice.toFixed(2));
}
// @ts-ignore
[order,];
const __VLS_export = (await import('vue')).defineComponent({});
export default {};
//# sourceMappingURL=OrderView.vue.js.map