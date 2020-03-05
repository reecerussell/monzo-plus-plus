import React from "react";
import { Container, Dropdown, Image, Menu, Icon } from "semantic-ui-react";
import { Link } from "react-router-dom";
import * as User from "../../utils/user";

const Header = () => {
	const adminMenu = User.IsInRole("Admin") ? (
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

				<Dropdown
					item
					simple
					text={
						<>
							<Icon name="user" />
							{User.GetUsername()}
						</>
					}
				>
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
