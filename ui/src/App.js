import React, { Suspense } from "react";
import "semantic-ui-css/semantic.min.css";
import Layout from "./components/layout/layout";
import {
	BrowserRouter,
	HashRouter as Router,
	Route,
	Switch,
} from "react-router-dom";
import { Loader } from "semantic-ui-react";

import routes from "./routes";

const Docs = React.lazy(() => import("./views/docs/index"));

const App = () => (
	<Suspense fallback={<Loader active={true} />}>
		<BrowserRouter>
			<Route path="/" exact>
				<Router>
					<Layout>
						<Switch>
							{routes.map((route, idx) => (
								<Route
									key={idx}
									name={route.name}
									exact={route.exact}
									path={route.path}
									render={props => (
										<route.component {...props} />
									)}
								/>
							))}
						</Switch>
					</Layout>
				</Router>
			</Route>
			<Route path="/docs">
				<Docs />
			</Route>
		</BrowserRouter>
	</Suspense>
);

export default App;
