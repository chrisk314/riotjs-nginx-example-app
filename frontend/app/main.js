import '@riotjs/hot-reload'
import {component, register} from 'riot'
import {Route, Router} from '@riotjs/route'

import App from './app.riot'
import AppMain from './app-main.riot'
import AppNavi from './app-navi.riot'

register('router', Router)
register('route', Route)

register('app-main', AppMain)
register('app-navi', AppNavi)

component(App)(document.getElementById('app'), {})
