import React from "react";
import {
	Container,
	Divider,
	Grid,
	Header,
	Image,
	List,
	Segment,
} from "semantic-ui-react";

const Footer = () => (
	<Segment
		inverted
		vertical
		style={{ margin: "5em 0em 0em", padding: "5em 0em" }}
	>
		<Container textAlign="center">
			<Grid divided inverted stackable>
				<Grid.Column width={3}>
					<Header inverted as="h4" content="Info" />
					<List link inverted>
						<List.Item as="a" href="/">
							Home
						</List.Item>
					</List>
				</Grid.Column>
				<Grid.Column width={3}>
					<Header inverted as="h4" content="Developer" />
					<List link inverted>
						<List.Item as="a" href="/docs">
							API Docs
						</List.Item>
						<List.Item
							as="a"
							href="https://github.com/reecerussell/monzo-plus-plus"
						>
							Source
						</List.Item>
						<List.Item as="a" href="https://reece-russell.co.uk">
							Reece Russell
						</List.Item>
					</List>
				</Grid.Column>
				<Grid.Column width={10}>
					<Header inverted as="h4" content="Footer Header" />
					<p>
						Extra space for a call to action inside the footer that
						could help re-engage users.
					</p>
				</Grid.Column>
			</Grid>

			<Divider inverted section />
			<Image centered size="mini" src="/logo.png" />
			<List horizontal inverted divided link size="small">
				<List.Item as="a" href="#">
					Site Map
				</List.Item>
				<List.Item as="a" href="#">
					Contact Us
				</List.Item>
				<List.Item as="a" href="#">
					Terms and Conditions
				</List.Item>
				<List.Item as="a" href="#">
					Privacy Policy
				</List.Item>
			</List>
		</Container>
	</Segment>
);

export default Footer;
