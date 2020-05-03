import '@riotjs/hot-reload'
import {component, register} from 'riot'
import {Route, Router} from '@riotjs/route'

import App from './shared/app.riot'
import AppMain from './shared/app-main.riot'
import AppNavi from './shared/app-navi.riot'

register('router', Router)
register('route', Route)

register('app-main', AppMain)
register('app-navi', AppNavi)

component(App)(document.getElementById('app'), {})
