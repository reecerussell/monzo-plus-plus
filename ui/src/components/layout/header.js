import React, { useState, useEffect } from "react";
import { Container, Dropdown, Image, Menu } from "semantic-ui-react";
import { Link } from "react-router-dom";
import { Authorise } from "../../containers/login";
import * as User from "../../utils/user";

const Header = () => {
	const [state, setState] = useState(0);

	const handleLogout = (e) => {
		e.preventDefault();

		User.Logout();
		window.location.reload();
	};

	useEffect(() => {
		User.SubscribeLogin("header", () => setState(state + 1));
		User.SubscribeLogout("header", () => setState(state + 1));

		return () => {
			User.UnsubscripeLogin("header");
			User.UnsubscribeLogout("header");
		};
	});

	return (
		<Menu fixed="top" inverted>
			<Container>
				<Menu.Item as="a" header href="/">
					<Image
						size="mini"
						src="https://lh3.googleusercontent.com/iDeb12CKMVdgDqBD9yJ9UehaWkKXFdPMtuUA8Jt0sOvxXzOm21qNGbA6D5_gdDZtAk4=w300"
						style={{ marginRight: "1.5em" }}
					/>
					Monzo++
				</Menu.Item>
				<Menu.Item as={Link} to="/">
					Home
				</Menu.Item>

				<Authorise contentOnly={true}>
					<Dropdown item simple text={User.GetUsername()}>
						<Dropdown.Menu>
							<Dropdown.Item>
								<span className="text">
									<Link to="/account">My Account</Link>
								</span>
							</Dropdown.Item>
							<Dropdown.Item>
								<span className="text">
									<a href="#" onClick={handleLogout}>
										Logout
									</a>
								</span>
							</Dropdown.Item>
							<Authorise
								contentOnly={true}
								roles={[
									"Admin",
									"Role Manager",
									"User Manager",
									"Plugin Manager",
								]}
							>
								<Dropdown.Divider />
								<Dropdown.Header>Admin</Dropdown.Header>
							</Authorise>
							<Authorise
								roles={["Admin", "User Manager"]}
								contentOnly={true}
							>
								<Dropdown.Item>
									<i className="dropdown icon" />
									<span className="text">
										<Link to="/users">Users</Link>
									</span>
									<Dropdown.Menu>
										<Dropdown.Item>
											<span className="text">
												<Link to="/users/pending">
													Pending
												</Link>
											</span>
										</Dropdown.Item>
									</Dropdown.Menu>
								</Dropdown.Item>
							</Authorise>
							<Authorise
								roles={["Admin", "Role Manager"]}
								contentOnly={true}
							>
								<Dropdown.Item>
									<span className="text">
										<Link to="/roles">Roles</Link>
									</span>
								</Dropdown.Item>
							</Authorise>
							<Authorise
								roles={["Admin", "Plugin Manager"]}
								contentOnly={true}
							>
								<Dropdown.Item>
									<span className="text">
										<Link to="/plugins">Plugins</Link>
									</span>
								</Dropdown.Item>
							</Authorise>
						</Dropdown.Menu>
					</Dropdown>
				</Authorise>
			</Container>
		</Menu>
	);
};

export default Header;
