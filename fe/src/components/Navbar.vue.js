/// <reference types="C:/Users/Ruslan/GolandProjects/diploma/fe/node_modules/@vue/language-core/types/template-helpers.d.ts" />
/// <reference types="C:/Users/Ruslan/GolandProjects/diploma/fe/node_modules/@vue/language-core/types/props-fallback.d.ts" />
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { ShoppingCart, User, Menu, X, Pizza } from 'lucide-vue-next';
import { useCartStore } from '../store/cart';
const router = useRouter();
const isMenuOpen = ref(false);
const cartStore = useCartStore();
const navLinks = [
    { name: 'Menu', path: '/' },
    { name: 'About', path: '/about' },
];
const toggleMenu = () => {
    isMenuOpen.value = !isMenuOpen.value;
};
const __VLS_ctx = {
    ...{},
    ...{},
};
let __VLS_components;
let __VLS_intrinsics;
let __VLS_directives;
__VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
    ...{ class: "navbar bg-primary text-primary-content shadow-lg sticky top-0 z-50" },
});
/** @type {__VLS_StyleScopedClasses['navbar']} */ ;
/** @type {__VLS_StyleScopedClasses['bg-primary']} */ ;
/** @type {__VLS_StyleScopedClasses['text-primary-content']} */ ;
/** @type {__VLS_StyleScopedClasses['shadow-lg']} */ ;
/** @type {__VLS_StyleScopedClasses['sticky']} */ ;
/** @type {__VLS_StyleScopedClasses['top-0']} */ ;
/** @type {__VLS_StyleScopedClasses['z-50']} */ ;
__VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
    ...{ class: "container mx-auto" },
});
/** @type {__VLS_StyleScopedClasses['container']} */ ;
/** @type {__VLS_StyleScopedClasses['mx-auto']} */ ;
__VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
    ...{ class: "flex-1" },
});
/** @type {__VLS_StyleScopedClasses['flex-1']} */ ;
let __VLS_0;
/** @ts-ignore @type {typeof __VLS_components.routerLink | typeof __VLS_components.RouterLink | typeof __VLS_components.routerLink | typeof __VLS_components.RouterLink} */
routerLink;
// @ts-ignore
const __VLS_1 = __VLS_asFunctionalComponent1(__VLS_0, new __VLS_0({
    to: "/",
    ...{ class: "btn btn-ghost normal-case text-xl flex items-center gap-2" },
}));
const __VLS_2 = __VLS_1({
    to: "/",
    ...{ class: "btn btn-ghost normal-case text-xl flex items-center gap-2" },
}, ...__VLS_functionalComponentArgsRest(__VLS_1));
/** @type {__VLS_StyleScopedClasses['btn']} */ ;
/** @type {__VLS_StyleScopedClasses['btn-ghost']} */ ;
/** @type {__VLS_StyleScopedClasses['normal-case']} */ ;
/** @type {__VLS_StyleScopedClasses['text-xl']} */ ;
/** @type {__VLS_StyleScopedClasses['flex']} */ ;
/** @type {__VLS_StyleScopedClasses['items-center']} */ ;
/** @type {__VLS_StyleScopedClasses['gap-2']} */ ;
const { default: __VLS_5 } = __VLS_3.slots;
let __VLS_6;
/** @ts-ignore @type {typeof __VLS_components.Pizza} */
Pizza;
// @ts-ignore
const __VLS_7 = __VLS_asFunctionalComponent1(__VLS_6, new __VLS_6({
    ...{ class: "w-8 h-8" },
}));
const __VLS_8 = __VLS_7({
    ...{ class: "w-8 h-8" },
}, ...__VLS_functionalComponentArgsRest(__VLS_7));
/** @type {__VLS_StyleScopedClasses['w-8']} */ ;
/** @type {__VLS_StyleScopedClasses['h-8']} */ ;
__VLS_asFunctionalElement1(__VLS_intrinsics.span, __VLS_intrinsics.span)({
    ...{ class: "font-bold tracking-tight" },
});
/** @type {__VLS_StyleScopedClasses['font-bold']} */ ;
/** @type {__VLS_StyleScopedClasses['tracking-tight']} */ ;
var __VLS_3;
__VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
    ...{ class: "hidden md:flex items-center gap-4 px-4" },
});
/** @type {__VLS_StyleScopedClasses['hidden']} */ ;
/** @type {__VLS_StyleScopedClasses['md:flex']} */ ;
/** @type {__VLS_StyleScopedClasses['items-center']} */ ;
/** @type {__VLS_StyleScopedClasses['gap-4']} */ ;
/** @type {__VLS_StyleScopedClasses['px-4']} */ ;
for (const [link] of __VLS_vFor((__VLS_ctx.navLinks))) {
    let __VLS_11;
    /** @ts-ignore @type {typeof __VLS_components.routerLink | typeof __VLS_components.RouterLink | typeof __VLS_components.routerLink | typeof __VLS_components.RouterLink} */
    routerLink;
    // @ts-ignore
    const __VLS_12 = __VLS_asFunctionalComponent1(__VLS_11, new __VLS_11({
        key: (link.path),
        to: (link.path),
        ...{ class: "btn btn-ghost" },
    }));
    const __VLS_13 = __VLS_12({
        key: (link.path),
        to: (link.path),
        ...{ class: "btn btn-ghost" },
    }, ...__VLS_functionalComponentArgsRest(__VLS_12));
    /** @type {__VLS_StyleScopedClasses['btn']} */ ;
    /** @type {__VLS_StyleScopedClasses['btn-ghost']} */ ;
    const { default: __VLS_16 } = __VLS_14.slots;
    (link.name);
    // @ts-ignore
    [navLinks,];
    var __VLS_14;
    // @ts-ignore
    [];
}
__VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
    ...{ class: "flex-none flex items-center gap-2" },
});
/** @type {__VLS_StyleScopedClasses['flex-none']} */ ;
/** @type {__VLS_StyleScopedClasses['flex']} */ ;
/** @type {__VLS_StyleScopedClasses['items-center']} */ ;
/** @type {__VLS_StyleScopedClasses['gap-2']} */ ;
let __VLS_17;
/** @ts-ignore @type {typeof __VLS_components.routerLink | typeof __VLS_components.RouterLink | typeof __VLS_components.routerLink | typeof __VLS_components.RouterLink} */
routerLink;
// @ts-ignore
const __VLS_18 = __VLS_asFunctionalComponent1(__VLS_17, new __VLS_17({
    to: "/cart",
    ...{ class: "btn btn-ghost btn-circle" },
}));
const __VLS_19 = __VLS_18({
    to: "/cart",
    ...{ class: "btn btn-ghost btn-circle" },
}, ...__VLS_functionalComponentArgsRest(__VLS_18));
/** @type {__VLS_StyleScopedClasses['btn']} */ ;
/** @type {__VLS_StyleScopedClasses['btn-ghost']} */ ;
/** @type {__VLS_StyleScopedClasses['btn-circle']} */ ;
const { default: __VLS_22 } = __VLS_20.slots;
__VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
    ...{ class: "indicator" },
});
/** @type {__VLS_StyleScopedClasses['indicator']} */ ;
let __VLS_23;
/** @ts-ignore @type {typeof __VLS_components.ShoppingCart} */
ShoppingCart;
// @ts-ignore
const __VLS_24 = __VLS_asFunctionalComponent1(__VLS_23, new __VLS_23({
    ...{ class: "h-5 w-5" },
}));
const __VLS_25 = __VLS_24({
    ...{ class: "h-5 w-5" },
}, ...__VLS_functionalComponentArgsRest(__VLS_24));
/** @type {__VLS_StyleScopedClasses['h-5']} */ ;
/** @type {__VLS_StyleScopedClasses['w-5']} */ ;
if (__VLS_ctx.cartStore.totalItems > 0) {
    __VLS_asFunctionalElement1(__VLS_intrinsics.span, __VLS_intrinsics.span)({
        ...{ class: "badge badge-sm indicator-item badge-secondary" },
    });
    /** @type {__VLS_StyleScopedClasses['badge']} */ ;
    /** @type {__VLS_StyleScopedClasses['badge-sm']} */ ;
    /** @type {__VLS_StyleScopedClasses['indicator-item']} */ ;
    /** @type {__VLS_StyleScopedClasses['badge-secondary']} */ ;
    (__VLS_ctx.cartStore.totalItems);
}
// @ts-ignore
[cartStore, cartStore,];
var __VLS_20;
__VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
    ...{ class: "dropdown dropdown-end" },
});
/** @type {__VLS_StyleScopedClasses['dropdown']} */ ;
/** @type {__VLS_StyleScopedClasses['dropdown-end']} */ ;
__VLS_asFunctionalElement1(__VLS_intrinsics.label, __VLS_intrinsics.label)({
    tabindex: "0",
    ...{ class: "btn btn-ghost btn-circle avatar" },
});
/** @type {__VLS_StyleScopedClasses['btn']} */ ;
/** @type {__VLS_StyleScopedClasses['btn-ghost']} */ ;
/** @type {__VLS_StyleScopedClasses['btn-circle']} */ ;
/** @type {__VLS_StyleScopedClasses['avatar']} */ ;
let __VLS_28;
/** @ts-ignore @type {typeof __VLS_components.User} */
User;
// @ts-ignore
const __VLS_29 = __VLS_asFunctionalComponent1(__VLS_28, new __VLS_28({
    ...{ class: "h-5 w-5" },
}));
const __VLS_30 = __VLS_29({
    ...{ class: "h-5 w-5" },
}, ...__VLS_functionalComponentArgsRest(__VLS_29));
/** @type {__VLS_StyleScopedClasses['h-5']} */ ;
/** @type {__VLS_StyleScopedClasses['w-5']} */ ;
__VLS_asFunctionalElement1(__VLS_intrinsics.ul, __VLS_intrinsics.ul)({
    tabindex: "0",
    ...{ class: "menu menu-sm dropdown-content mt-3 z-[1] p-2 shadow bg-base-100 text-base-content rounded-box w-52" },
});
/** @type {__VLS_StyleScopedClasses['menu']} */ ;
/** @type {__VLS_StyleScopedClasses['menu-sm']} */ ;
/** @type {__VLS_StyleScopedClasses['dropdown-content']} */ ;
/** @type {__VLS_StyleScopedClasses['mt-3']} */ ;
/** @type {__VLS_StyleScopedClasses['z-[1]']} */ ;
/** @type {__VLS_StyleScopedClasses['p-2']} */ ;
/** @type {__VLS_StyleScopedClasses['shadow']} */ ;
/** @type {__VLS_StyleScopedClasses['bg-base-100']} */ ;
/** @type {__VLS_StyleScopedClasses['text-base-content']} */ ;
/** @type {__VLS_StyleScopedClasses['rounded-box']} */ ;
/** @type {__VLS_StyleScopedClasses['w-52']} */ ;
__VLS_asFunctionalElement1(__VLS_intrinsics.li, __VLS_intrinsics.li)({});
let __VLS_33;
/** @ts-ignore @type {typeof __VLS_components.routerLink | typeof __VLS_components.RouterLink | typeof __VLS_components.routerLink | typeof __VLS_components.RouterLink} */
routerLink;
// @ts-ignore
const __VLS_34 = __VLS_asFunctionalComponent1(__VLS_33, new __VLS_33({
    to: "/login",
}));
const __VLS_35 = __VLS_34({
    to: "/login",
}, ...__VLS_functionalComponentArgsRest(__VLS_34));
const { default: __VLS_38 } = __VLS_36.slots;
// @ts-ignore
[];
var __VLS_36;
__VLS_asFunctionalElement1(__VLS_intrinsics.li, __VLS_intrinsics.li)({});
let __VLS_39;
/** @ts-ignore @type {typeof __VLS_components.routerLink | typeof __VLS_components.RouterLink | typeof __VLS_components.routerLink | typeof __VLS_components.RouterLink} */
routerLink;
// @ts-ignore
const __VLS_40 = __VLS_asFunctionalComponent1(__VLS_39, new __VLS_39({
    to: "/register",
}));
const __VLS_41 = __VLS_40({
    to: "/register",
}, ...__VLS_functionalComponentArgsRest(__VLS_40));
const { default: __VLS_44 } = __VLS_42.slots;
// @ts-ignore
[];
var __VLS_42;
__VLS_asFunctionalElement1(__VLS_intrinsics.li, __VLS_intrinsics.li)({});
__VLS_asFunctionalElement1(__VLS_intrinsics.a, __VLS_intrinsics.a)({});
__VLS_asFunctionalElement1(__VLS_intrinsics.li, __VLS_intrinsics.li)({});
__VLS_asFunctionalElement1(__VLS_intrinsics.a, __VLS_intrinsics.a)({});
__VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
    ...{ class: "divider my-0" },
});
/** @type {__VLS_StyleScopedClasses['divider']} */ ;
/** @type {__VLS_StyleScopedClasses['my-0']} */ ;
__VLS_asFunctionalElement1(__VLS_intrinsics.li, __VLS_intrinsics.li)({});
__VLS_asFunctionalElement1(__VLS_intrinsics.a, __VLS_intrinsics.a)({});
__VLS_asFunctionalElement1(__VLS_intrinsics.button, __VLS_intrinsics.button)({
    ...{ onClick: (__VLS_ctx.toggleMenu) },
    ...{ class: "btn btn-ghost btn-circle md:hidden" },
});
/** @type {__VLS_StyleScopedClasses['btn']} */ ;
/** @type {__VLS_StyleScopedClasses['btn-ghost']} */ ;
/** @type {__VLS_StyleScopedClasses['btn-circle']} */ ;
/** @type {__VLS_StyleScopedClasses['md:hidden']} */ ;
if (!__VLS_ctx.isMenuOpen) {
    let __VLS_45;
    /** @ts-ignore @type {typeof __VLS_components.Menu} */
    Menu;
    // @ts-ignore
    const __VLS_46 = __VLS_asFunctionalComponent1(__VLS_45, new __VLS_45({
        ...{ class: "h-5 w-5" },
    }));
    const __VLS_47 = __VLS_46({
        ...{ class: "h-5 w-5" },
    }, ...__VLS_functionalComponentArgsRest(__VLS_46));
    /** @type {__VLS_StyleScopedClasses['h-5']} */ ;
    /** @type {__VLS_StyleScopedClasses['w-5']} */ ;
}
else {
    let __VLS_50;
    /** @ts-ignore @type {typeof __VLS_components.X} */
    X;
    // @ts-ignore
    const __VLS_51 = __VLS_asFunctionalComponent1(__VLS_50, new __VLS_50({
        ...{ class: "h-5 w-5" },
    }));
    const __VLS_52 = __VLS_51({
        ...{ class: "h-5 w-5" },
    }, ...__VLS_functionalComponentArgsRest(__VLS_51));
    /** @type {__VLS_StyleScopedClasses['h-5']} */ ;
    /** @type {__VLS_StyleScopedClasses['w-5']} */ ;
}
if (__VLS_ctx.isMenuOpen) {
    __VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
        ...{ class: "md:hidden bg-base-200 p-4 border-b border-base-300" },
    });
    /** @type {__VLS_StyleScopedClasses['md:hidden']} */ ;
    /** @type {__VLS_StyleScopedClasses['bg-base-200']} */ ;
    /** @type {__VLS_StyleScopedClasses['p-4']} */ ;
    /** @type {__VLS_StyleScopedClasses['border-b']} */ ;
    /** @type {__VLS_StyleScopedClasses['border-base-300']} */ ;
    __VLS_asFunctionalElement1(__VLS_intrinsics.ul, __VLS_intrinsics.ul)({
        ...{ class: "menu w-full gap-2" },
    });
    /** @type {__VLS_StyleScopedClasses['menu']} */ ;
    /** @type {__VLS_StyleScopedClasses['w-full']} */ ;
    /** @type {__VLS_StyleScopedClasses['gap-2']} */ ;
    for (const [link] of __VLS_vFor((__VLS_ctx.navLinks))) {
        __VLS_asFunctionalElement1(__VLS_intrinsics.li, __VLS_intrinsics.li)({
            key: (link.path),
        });
        let __VLS_55;
        /** @ts-ignore @type {typeof __VLS_components.routerLink | typeof __VLS_components.RouterLink | typeof __VLS_components.routerLink | typeof __VLS_components.RouterLink} */
        routerLink;
        // @ts-ignore
        const __VLS_56 = __VLS_asFunctionalComponent1(__VLS_55, new __VLS_55({
            ...{ 'onClick': {} },
            to: (link.path),
        }));
        const __VLS_57 = __VLS_56({
            ...{ 'onClick': {} },
            to: (link.path),
        }, ...__VLS_functionalComponentArgsRest(__VLS_56));
        let __VLS_60;
        const __VLS_61 = ({ click: {} },
            { onClick: (...[$event]) => {
                    if (!(__VLS_ctx.isMenuOpen))
                        return;
                    __VLS_ctx.isMenuOpen = false;
                    // @ts-ignore
                    [navLinks, toggleMenu, isMenuOpen, isMenuOpen, isMenuOpen,];
                } });
        const { default: __VLS_62 } = __VLS_58.slots;
        (link.name);
        // @ts-ignore
        [];
        var __VLS_58;
        var __VLS_59;
        // @ts-ignore
        [];
    }
}
// @ts-ignore
[];
const __VLS_export = (await import('vue')).defineComponent({});
export default {};
//# sourceMappingURL=Navbar.vue.js.map