import React, { Suspense } from "react";
import "semantic-ui-css/semantic.min.css";
import Layout from "./components/layout/layout";
import { HashRouter as Router, Route, Switch } from "react-router-dom";
import { Loader } from "semantic-ui-react";

import routes from "./routes";

const App = () => (
	<Suspense fallback={<Loader active={true} />}>
		<Router>
			<Layout>
				<Switch>
					{routes.map((route, idx) => (
						<Route
							key={idx}
							name={route.name}
							exact={route.exact}
							path={route.path}
							render={props => <route.component {...props} />}
						/>
					))}
				</Switch>
			</Layout>
		</Router>
	</Suspense>
);

export default App;
