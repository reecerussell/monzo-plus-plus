import React, { useState, useEffect, createRef } from "react";
import { Menu, Responsive, Sticky, Ref, Grid } from "semantic-ui-react";
import { Redirect, useLocation } from "react-router-dom";

const Layout = ({ children }) => {
	const { pathname } = useLocation();
	const [activePath, setPath] = useState(pathname);
	const [redirect, setRedirect] = useState(null);
	const [menuWidth, setMenuWidth] = useState(16);
	const contextRef = createRef();

	const handleResponsiveUpdate = () => {
		const width = window.innerWidth;

		if (
			width <= Responsive.onlyWidescreen.maxWidth &&
			width >= Responsive.onlyWidescreen.minWidth
		) {
			setMenuWidth(3);
		}
		if (
			width <= Responsive.onlyLargeScreen.maxWidth &&
			width >= Responsive.onlyLargeScreen.minWidth
		) {
			setMenuWidth(4);
		}
		if (
			width <= Responsive.onlyComputer.maxWidth &&
			width >= Responsive.onlyComputer.minWidth
		) {
			setMenuWidth(5);
		}
		if (
			width <= Responsive.onlyTablet.maxWidth &&
			width >= Responsive.onlyTablet.minWidth
		) {
			setMenuWidth(5);
		}
		if (
			width <= Responsive.onlyMobile.maxWidth &&
			width >= Responsive.onlyMobile.minWidth
		) {
			setMenuWidth(16);
		}
	};

	useEffect(() => {
		if (activePath === pathname) {
			setRedirect(null);
		} else {
			setRedirect(activePath);
		}
	}, [activePath]);

	useEffect(() => {
		handleResponsiveUpdate();
	}, []);

	if (redirect) {
		return <Redirect to={redirect} />;
	}

	return (
		<Responsive onUpdate={handleResponsiveUpdate}>
			<Grid>
				<Grid.Row>
					<Grid.Column width={menuWidth}>
						<Ref innerRef={contextRef}>
							<Sticky context={contextRef}>
								<Menu pointing vertical fluid>
									<Menu.Item
										name={"Details"}
										active={activePath === "/account"}
										onClick={() => setPath("/account")}
									/>
									<Menu.Item
										name={"Change Password"}
										active={
											activePath ===
											"/account/changepassword"
										}
										onClick={() =>
											setPath("/account/changepassword")
										}
									/>
									<Menu.Item
										name={"Plugins"}
										active={
											activePath === "/account/plugins"
										}
										onClick={() =>
											setPath("/account/plugins")
										}
									/>
									<Menu.Item
										name={"Delete"}
										color="red"
										active={
											activePath === "/account/delete"
										}
										onClick={() =>
											setPath("/account/delete")
										}
									/>
								</Menu>
							</Sticky>
						</Ref>
					</Grid.Column>
					<Grid.Column width={menuWidth === 16 ? 16 : 16 - menuWidth}>
						{children}
					</Grid.Column>
				</Grid.Row>
			</Grid>
		</Responsive>
	);
};

export default Layout;
