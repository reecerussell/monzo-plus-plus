import React, { useState, useEffect, createRef } from "react";
import { Menu, Sticky, Ref, Grid, Header, Divider } from "semantic-ui-react";
import { Redirect, useLocation } from "react-router-dom";
import SetAccountContainer from "../../containers/account/setAccount";
import * as User from "../../utils/user";

const Layout = ({ children }) => {
	const { pathname } = useLocation();
	const [activePath, setPath] = useState(pathname);
	const [redirect, setRedirect] = useState(null);
	const contextRef = createRef();

	useEffect(() => {
		if (activePath === pathname) {
			setRedirect(null);
		} else {
			setRedirect(activePath);
		}
	}, [activePath, pathname]);

	if (redirect) {
		return <Redirect to={redirect} />;
	}

	return (
		<>
			<Header as="h1">Manage your account</Header>
			<Divider />
			<Grid stackable>
				<Grid.Row>
					<Grid.Column width={4}>
						<Ref innerRef={contextRef}>
							<Sticky context={contextRef}>
								<Menu vertical fluid>
									{User.HasAccount() ? (
										<>
											<Menu.Item
												name={"Details"}
												active={
													activePath === "/account"
												}
												onClick={() =>
													setPath("/account")
												}
											/>
											<Menu.Item
												name={"Change Password"}
												active={
													activePath ===
													"/account/changepassword"
												}
												onClick={() =>
													setPath(
														"/account/changepassword"
													)
												}
											/>
											<Menu.Item
												name={"Plugins"}
												active={
													activePath ===
													"/account/plugins"
												}
												onClick={() =>
													setPath("/account/plugins")
												}
											/>
											<Menu.Item
												name={"Delete"}
												color="red"
												active={
													activePath ===
													"/account/delete"
												}
												onClick={() =>
													setPath("/account/delete")
												}
											/>
										</>
									) : (
										<Menu.Item
											name="Select an account"
											active={true}
										/>
									)}
								</Menu>
							</Sticky>
						</Ref>
					</Grid.Column>
					<Grid.Column width={12}>
						{User.HasAccount() ? children : <SetAccountContainer />}
					</Grid.Column>
				</Grid.Row>
			</Grid>
		</>
	);
};

export default Layout;
