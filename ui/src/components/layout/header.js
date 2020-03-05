import React from "react";
import { Container, Dropdown, Image, Menu } from "semantic-ui-react";
import * as User from "../../utils/user";

const Header = () => {
	const adminMenu = User.IsInRole("Admin") ? (
		<>
			<Dropdown.Divider />
			<Dropdown.Header>Admin</Dropdown.Header>
			<Dropdown.Item>
				<i className="dropdown icon" />
				<span className="text">Users</span>
				<Dropdown.Menu>
					<Dropdown.Item>All</Dropdown.Item>
					<Dropdown.Item>Pending</Dropdown.Item>
				</Dropdown.Menu>
			</Dropdown.Item>
			<Dropdown.Item>Roles</Dropdown.Item>{" "}
		</>
	) : null;

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

				<Dropdown item simple text="Dropdown">
					<Dropdown.Menu>
						<Dropdown.Item>My Account</Dropdown.Item>
						<Dropdown.Item>Logout</Dropdown.Item>
						{adminMenu}
					</Dropdown.Menu>
				</Dropdown>
			</Container>
		</Menu>
	);
};

export default Header;
