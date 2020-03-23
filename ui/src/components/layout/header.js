import React, { useState, useEffect } from "react";
import { Container, Dropdown, Image, Menu } from "semantic-ui-react";
import { Link } from "react-router-dom";
import * as User from "../../utils/user";

const AdminMenu = () => {
	if (!User.IsInRole("Admin")) {
		return null;
	}

	return (
		<>
			<Dropdown.Divider />
			<Dropdown.Header>Admin</Dropdown.Header>
			<Dropdown.Item>
				<i className="dropdown icon" />
				<span className="text">
					<Link to="/users">Users</Link>
				</span>
				<Dropdown.Menu>
					<Dropdown.Item>
						<span className="text">
							<Link to="/users/pending">Pending</Link>
						</span>
					</Dropdown.Item>
				</Dropdown.Menu>
			</Dropdown.Item>
			<Dropdown.Item>
				<span className="text">
					<Link to="/roles">Roles</Link>
				</span>
			</Dropdown.Item>
		</>
	);
};

const UserMenu = () => {
	if (!User.IsAuthenticated()) {
		return (
			<Menu.Item as={Link} to="/login">
				Login
			</Menu.Item>
		);
	}

	const handleLogout = e => {
		e.preventDefault();

		User.Logout();
		window.location.reload();
	};

	return (
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
				<AdminMenu />
			</Dropdown.Menu>
		</Dropdown>
	);
};

const Header = () => {
	const [state, setState] = useState(1);

	useEffect(() => {
		User.SubscribeLogin("header", () => setState(1));
		User.SubscribeLogout("header", () => setState(null));

		return () => {
			User.UnsubscripeLogin("header");
			User.UnsubscribeLogout("header");
		};
	});

	return (
		<Menu fixed="top" inverted>
			<Container>
				<Menu.Item as="a" header>
					<Image
						size="mini"
						src="https://lh3.googleusercontent.com/iDeb12CKMVdgDqBD9yJ9UehaWkKXFdPMtuUA8Jt0sOvxXzOm21qNGbA6D5_gdDZtAk4=w300"
						style={{ marginRight: "1.5em" }}
					/>
					Monzo++
				</Menu.Item>
				<Menu.Item as="a">Home</Menu.Item>

				<UserMenu />
			</Container>
		</Menu>
	);
};

export default Header;
