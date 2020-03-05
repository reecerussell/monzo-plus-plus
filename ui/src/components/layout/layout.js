import React from "react";
import { Container } from "semantic-ui-react";
import Footer from "./footer";
import Header from "./header";

const Layout = ({ children }) => (
	<>
		<Header />

		<Container style={{ marginTop: "7em" }}>{children}</Container>

		<Footer />
	</>
);

export default Layout;
