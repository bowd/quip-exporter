import Navigo from "navigo"

const router = new Navigo(null, true, '#!');
router.on('/', () => {
    console.log('root')
}).resolve();