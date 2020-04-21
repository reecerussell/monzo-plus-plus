import React from "react";
import {
	Item,
	Form,
	Loader,
	Button,
	Message,
	Icon,
	Label,
} from "semantic-ui-react";
import { Link } from "react-router-dom";

const List = ({
	plugins,
	error,
	loading,
	searchTerm,
	handleUpdateSearch,
	handleSearch,
}) => {
	const renderRow = (plugin, idx) => (
		<Item key={idx}>
			<Item.Content>
				<Item.Header as="h4">
					<Link to={"/plugins/edit/" + plugin.id}>
						{plugin.displayName}
					</Link>
				</Item.Header>
				<Item.Meta>{plugin.name}</Item.Meta>
				<Item.Description>{plugin.description}</Item.Description>
				<Item.Extra>
					<Label>
						<Icon name="user"></Icon>
						{plugin.consumedBy}
					</Label>
					<Button
						as={Link}
						to={"/plugins/edit/" + plugin.id}
						color="blue"
						size="small"
						floated="right"
					>
						View <Icon name="arrow right" />
					</Button>
				</Item.Extra>
			</Item.Content>
		</Item>
	);

	const errorItem =
		plugins.length < 1 ? (
			searchTerm.length < 1 ? (
				<Item>
					<Item.Content>
						<Item.Header>No Plugins Found</Item.Header>
						<Item.Description>
							No plugins currently exist. If you create one, it
							will appear here.
						</Item.Description>
					</Item.Content>
				</Item>
			) : (
				<Item>
					<Item.Content>
						<Item.Header>No Plugins Found</Item.Header>
						<Item.Description>
							No plugins were found matching your search. Try
							changing your query and trying again.
						</Item.Description>
					</Item.Content>
				</Item>
			)
		) : null;

	return (
		<>
			<Loader active={loading} />

			<Form onSubmit={handleSearch} error={error !== null}>
				<Message error header="An error occured!" content={error} />

				<Form.Group inline>
					<Form.Field>
						<input
							type="search"
							value={searchTerm}
							onChange={handleUpdateSearch}
							placeholder="Search..."
						/>
					</Form.Field>
					<Form.Field>
						<Button color="grey" type="submit">
							<Icon name="search" />
							Search
						</Button>
					</Form.Field>
				</Form.Group>
			</Form>

			<Item.Group divided>
				{errorItem}

				{plugins.map(renderRow)}
			</Item.Group>
		</>
	);
};

export default List;
