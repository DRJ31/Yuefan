const Foo = {template: '<h1>Foo</h1>'};
const Bar = {template: '<h1>Bar</h1>'};

const routes = [
    {path: '/foo', component: Foo},
    {path: '/bar', component: Bar}
];

const router = new VueRouter({
    routes
});

const app = new Vue({
    router,
    data: {
        restaurants: [
            '卫师傅',
            '麻辣烫'
        ]
    },
    methods: {
        add_restaurant: function () {
            if (this.$refs.restaurant_add.value !== "") {
                this.restaurants.push(this.$refs.restaurant_add.value);
                this.$refs.restaurant_add.value = "";
            }
        }
    }
}).$mount("#yuefan");