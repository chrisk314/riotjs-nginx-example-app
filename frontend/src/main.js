import '@riotjs/hot-reload'
import {component, register} from 'riot'
import {Route, Router} from '@riotjs/route'

import App from './shared/app.riot'

register('router', Router)
register('route', Route)

const mountApp = component(App)
const app = mountApp(document.getElementById('app'), {})
