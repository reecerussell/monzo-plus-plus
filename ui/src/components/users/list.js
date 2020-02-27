import React from "react";
import { Icon, Table, Message, Loader, Segment, Form } from "semantic-ui-react";
import { Link } from "react-router-dom";
import PropTypes from "prop-types";

const propTypes = {
	users: PropTypes.array,
	loading: PropTypes.bool.isRequired,
	error: PropTypes.string,
	searchTerm: PropTypes.string,
	updateSearchTerm: PropTypes.func,
	onSearch: PropTypes.func,
};
const defaultProps = {
	users: [],
	error: null,
};

const List = ({
	users,
	loading,
	error,
	searchTerm,
	updateSearchTerm,
	onSearch,
}) => {
	const errorMessage = error ? (
		<Message error header="An error occured!" content={error} />
	) : null;

	const getUsername = username => {
		if (searchTerm.length < 1) {
			return username;
		}

		const regex = new RegExp(searchTerm, "gi");

		return username.replace(
			regex,
			str => "<span style='background-color: yellow;'>" + str + "</span>"
		);
	};

	const rows =
		users && users.length > 0 ? (
			users.map((user, idx) => (
				<Table.Row key={idx}>
					<Table.Cell>
						<span
							dangerouslySetInnerHTML={{
								__html: getUsername(user.username),
							}}
						/>
					</Table.Cell>
					<Table.Cell>
						{user.enabled ? (
							<Icon name="check" />
						) : (
							<Icon name="close" />
						)}
					</Table.Cell>
					<Table.Cell>
						<Link to={"/users/details/" + user.id}>View</Link>
					</Table.Cell>
				</Table.Row>
			))
		) : (
			<Table.Row>
				<Table.Cell colSpan="3">
					<b>
						{searchTerm.length > 0
							? "No users were found for this search, try changing your query."
							: "No users exist."}
					</b>
				</Table.Cell>
			</Table.Row>
		);

	return (
		<>
			{errorMessage}
			<Segment>
				<Loader active={loading} />

				<Form onSubmit={onSearch}>
					<Form.Group>
						<Form.Input
							placeholder="Search..."
							value={searchTerm}
							onChange={e => updateSearchTerm(e.target.value)}
						/>
						<Form.Button content="Search" />
					</Form.Group>
				</Form>

				<Table striped>
					<Table.Header>
						<Table.Row>
							<Table.HeaderCell>Username</Table.HeaderCell>
							<Table.HeaderCell>Enabled</Table.HeaderCell>
							<Table.HeaderCell></Table.HeaderCell>
						</Table.Row>
					</Table.Header>

					<Table.Body>{rows}</Table.Body>
				</Table>
			</Segment>
		</>
	);
};

List.propTypes = propTypes;
List.defaultProps = defaultProps;

export default List;
