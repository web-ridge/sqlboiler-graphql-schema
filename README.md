# sqlboiler-graphql-schema

This program generates a grapql schema based on the generated sqlboiler structs we do this because we want to support the sqlboiler aliasing in our schema. Generating the schema is a good way too add type safety to queries and filters and prevent too much manual typing.

You can edit your schema like you want later and re-generate if your database changes. This program will create a merge conflict with your existing schema so you can choose to accept/reject changes.

## How to run

`go run github.com/web-ridge/sqlboiler-graphql-schema`

## Before running

- Install prettier globally (https://prettier.io/ `yarn global add prettier`)
- Install git command line (required to do three way merging)

## Other related projects from webRidge
- https://github.com/web-ridge/gqlgen-sqlboiler (Generates converts between your qqlgen scheme and sqlboiler, and in the future will generate basic resolvers)
- https://github.com/web-ridge/graphql-schema-react-native-app (Generated React Native (Web) app based on your GraphQL scheme, not finished yet.)

## Options

```
USAGE:
    [global options] command [command options] [arguments...]

COMMANDS:
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --input value   directory where the sqlboiler models are (default: "models")
   --output value  filepath for schema (default: "schema.graphql")
   --mutations     generate mutations for models
   --batch-update  generate batch update for models
   --batch-create  generate batch create for models
   --batch-delete  generate batch delete for models
   --help, -h      show help
   --version, -v   print the version
```

- Models
- Mutations (Followed best practices https://blog.apollographql.com/designing-graphql-mutations-e09de826ed97)

## Features

- Support for manual updating the schema and re-generating (doing a three way merge https://github.com/charlesvdv/go-three-way-merge)
- Generating basic models (100%)
- Generating basic queries (100%)
- Generating filter for array queries (100%)
- Generating filters for relationships (100%)
- Generating input for mutations (100%)
- Generating payload for mutations (100%)
- Generating mutations (100%)
- Generating mutations for array models (0% WIP)
- Generating pagination for array models (0% WIP)

## Future roadmap

- Tests / snapshots
- Edges / connections
- Detecting when relationship is many to many
- Adding node from to many-to-many relationships
- Removing node from many-to-many relationships

Feel free to help and make a PR.

```
go run github.com/webridge-git/sqlboiler-graphql-schema
```

## Example result

```graphql
type Address {
  id: ID!
  street: String
  houseNumber: String
  zipAddress: String
  city: String
  longitude: Float!
  latitude: Float!
  description: String
  name: String
  permission: Boolean
  addressStatus: AddressStatus!
  company: Company!
  contactPerson: Person!
  houseType: HouseType!
  owner: Person!
  userOrganization: UserOrganization!
  calamities: [Calamity]!
  people: [Person]!
  createdAt: Int
  deletedAt: Int
  updatedAt: Int
}

type AddressStatus {
  id: ID!
  name: String
  icon: String
  description: String
  order: Int!
  color: String!
  addresses: [Address]!
  createdAt: Int
  updatedAt: Int
  deletedAt: Int
}

type Calamity {
  id: ID!
  description: String
  priority: Int!
  totalBuildings: String
  street: String
  houseNumbers: String
  zipAddress: String
  longitude: String
  latitude: String
  city: String
  underground: String
  part: String
  color: String
  whatColor: String
  resolvedM2: Float
  resolvedDate: Int
  resolvedCompanyDescription: String
  resolvedUserDescription: String
  m2: Float
  amount: Int
  address: Address!
  calamityType: CalamityType!
  company: Company!
  houseType: HouseType!
  parent: Calamity!
  status: Status!
  user: User!
  userOrganization: UserOrganization!
  parentCalamities: [Calamity]!
  calamityAttributes: [CalamityAttribute]!
  calamityTags: [CalamityTag]!
  createdAt: Int
  updatedAt: Int
  deletedAt: Int
}

type CalamityAttribute {
  id: ID!
  key: String!
  type: String!
  value: String
  calamity: Calamity!
  deletedAt: Int
  createdAt: Int
  updatedAt: Int
}

type CalamityPicture {
  id: ID!
  name: String
  url: String!
  calamityId: ID!
  deletedAt: Int
  createdAt: Int
  updatedAt: Int
}

type CalamityTag {
  id: ID!
  calamity: Calamity!
  tag: Tag!
}

type CalamityType {
  id: ID!
  name: String!
  order: Int!
  image: String
  batchCreate: Boolean
  description: String
  defaultM2: Float!
  userOrganization: UserOrganization!
  calamities: [Calamity]!
  updatedAt: Int
  createdAt: Int
  deletedAt: Int
}

type Company {
  id: ID!
  name: String!
  backgroundColor: String!
  textColor: String!
  order: Int
  userOrganization: UserOrganization!
  addresses: [Address]!
  calamities: [Calamity]!
  companyUsers: [CompanyUser]!
  updatedAt: Int
  createdAt: Int
  deletedAt: Int
}

type CompanyUser {
  id: ID!
  company: Company!
  user: User!
  deletedAt: Int
  updatedAt: Int
  createdAt: Int
}

type FailedJob {
  id: ID!
  connection: String!
  queue: String!
  payload: String!
  exception: String!
  failedAt: Int!
}

type HouseType {
  id: ID!
  name: String!
  order: Int!
  image: String
  description: String
  userOrganization: UserOrganization!
  addresses: [Address]!
  calamities: [Calamity]!
  deletedAt: Int
  updatedAt: Int
  createdAt: Int
}

type Migration {
  id: ID!
  migration: String!
  batch: Int!
}

type Person {
  id: ID!
  firstName: String
  lastName: String
  description: String
  telephoneNumber: String
  email: String
  hasGivenPermission: Boolean
  name: String
  companyName: String
  address: Address!
  userOrganization: UserOrganization!
  contactPersonAddresses: [Address]!
  ownerAddresses: [Address]!
  deletedAt: Int
  createdAt: Int
  updatedAt: Int
}

type Role {
  id: ID!
  name: String!
  icon: String
  description: String
  users: [User]!
  updatedAt: Int
  deletedAt: Int
  createdAt: Int
}

type Status {
  id: ID!
  order: Int!
  name: String!
  icon: String
  color: String
  iconThumb: String
  calamities: [Calamity]!
  updatedAt: Int
  createdAt: Int
  deletedAt: Int
}

type Tag {
  id: ID!
  name: String!
  image: String
  userOrganization: UserOrganization!
  calamityTags: [CalamityTag]!
  createdAt: Int
  updatedAt: Int
  deletedAt: Int
}

type User {
  id: ID!
  name: String!
  lastName: String!
  email: String!
  password: String!
  rememberToken: String
  sendNotificationsOnNewCalamity: Boolean!
  role: Role!
  calamities: [Calamity]!
  companyUsers: [CompanyUser]!
  userUserOrganizations: [UserUserOrganization]!
  createdAt: Int
  deletedAt: Int
  updatedAt: Int
}

type UserOrganization {
  id: ID!
  name: String
  logo: String
  street: String
  city: String
  houseNumber: String
  zipCode: String
  telephoneNumber: String
  email: String
  primaryColor: String
  accentColor: String
  enablePermissionOnAddress: Boolean!
  enableContact: Boolean!
  enableOwner: Boolean!
  enableContactPerson: Boolean!
  newCalamityCompanyScreen: Boolean!
  newCalamityPriorityScreen: Boolean!
  newCalamityM2Screen: Boolean!
  addresses: [Address]!
  calamities: [Calamity]!
  calamityTypes: [CalamityType]!
  companies: [Company]!
  houseTypes: [HouseType]!
  people: [Person]!
  tags: [Tag]!
  userUserOrganizations: [UserUserOrganization]!
  createdAt: Int
  updatedAt: Int
  deletedAt: Int
}

type UserUserOrganization {
  id: ID!
  user: User!
  userOrganization: UserOrganization!
}

type FilterID {
  in: [ID!]
  notIn: [ID!]
}
type FilterString {
  equalTo: String
  in: [String!]
  notIn: [String!]

  startsWith: String
  endsWith: String
  contains: String

  startsWithStrict: String # Camel sensitive
  endsWithStrict: String # Camel sensitive
  containsStrict: String # Camel sensitive
}
type FilterInteger {
  equalTo: Int
  lessThan: Int
  lessThanOrEqualTo: Int
  moreThan: Int
  moreThanOrEqualTo: Int
  in: [Int!]
  notIn: [Int!]
}
type FilterBoolean {
  equalTo: Boolean
}

type AddressFilter {
  search: String
  where: AddressWhere
}

type AddressWhere {
  id: IDFilter
  street: StringFilter
  houseNumber: StringFilter
  zipAddress: StringFilter
  city: StringFilter
  longitude: FloatFilter
  latitude: FloatFilter
  ownerId: IDFilter
  contactPersonId: IDFilter
  userOrganizationId: IDFilter
  companyId: IDFilter
  houseTypeId: IDFilter
  description: StringFilter
  name: StringFilter
  permission: BooleanFilter
  addressStatusId: IDFilter
  addressStatus: AddressStatusWhere
  company: CompanyWhere
  contactPerson: PersonWhere
  houseType: HouseTypeWhere
  owner: PersonWhere
  userOrganization: UserOrganizationWhere
  calamities: CalamityWhere
  people: PersonWhere
  createdAt: IntFilter
  deletedAt: IntFilter
  updatedAt: IntFilter
  or: AddressWhere
  and: AddressWhere
}

type AddressStatusFilter {
  search: String
  where: AddressStatusWhere
}

type AddressStatusWhere {
  id: IDFilter
  name: StringFilter
  icon: StringFilter
  description: StringFilter
  order: IntFilter
  color: StringFilter
  addresses: AddressWhere
  createdAt: IntFilter
  updatedAt: IntFilter
  deletedAt: IntFilter
  or: AddressStatusWhere
  and: AddressStatusWhere
}

type CalamityFilter {
  search: String
  where: CalamityWhere
}

type CalamityWhere {
  id: IDFilter
  parentId: IDFilter
  description: StringFilter
  priority: IntFilter
  totalBuildings: StringFilter
  street: StringFilter
  houseNumbers: StringFilter
  zipAddress: StringFilter
  longitude: StringFilter
  latitude: StringFilter
  statusId: IDFilter
  companyId: IDFilter
  userId: IDFilter
  houseTypeId: IDFilter
  calamityTypeId: IDFilter
  city: StringFilter
  userOrganizationId: IDFilter
  addressId: IDFilter
  underground: StringFilter
  part: StringFilter
  color: StringFilter
  whatColor: StringFilter
  resolvedM2: FloatFilter
  resolvedDate: IntFilter
  resolvedCompanyDescription: StringFilter
  resolvedUserDescription: StringFilter
  m2: FloatFilter
  amount: IntFilter
  address: AddressWhere
  calamityType: CalamityTypeWhere
  company: CompanyWhere
  houseType: HouseTypeWhere
  parent: CalamityWhere
  status: StatusWhere
  user: UserWhere
  userOrganization: UserOrganizationWhere
  parentCalamities: CalamityWhere
  calamityAttributes: CalamityAttributeWhere
  calamityTags: CalamityTagWhere
  createdAt: IntFilter
  updatedAt: IntFilter
  deletedAt: IntFilter
  or: CalamityWhere
  and: CalamityWhere
}

type CalamityAttributeFilter {
  search: String
  where: CalamityAttributeWhere
}

type CalamityAttributeWhere {
  id: IDFilter
  key: StringFilter
  type: StringFilter
  value: StringFilter
  calamityId: IDFilter
  calamity: CalamityWhere
  deletedAt: IntFilter
  createdAt: IntFilter
  updatedAt: IntFilter
  or: CalamityAttributeWhere
  and: CalamityAttributeWhere
}

type CalamityPictureFilter {
  search: String
  where: CalamityPictureWhere
}

type CalamityPictureWhere {
  id: IDFilter
  name: StringFilter
  url: StringFilter
  calamityId: IDFilter
  deletedAt: IntFilter
  createdAt: IntFilter
  updatedAt: IntFilter
  or: CalamityPictureWhere
  and: CalamityPictureWhere
}

type CalamityTagFilter {
  search: String
  where: CalamityTagWhere
}

type CalamityTagWhere {
  id: IDFilter
  calamityId: IDFilter
  tagId: IDFilter
  calamity: CalamityWhere
  tag: TagWhere
  or: CalamityTagWhere
  and: CalamityTagWhere
}

type CalamityTypeFilter {
  search: String
  where: CalamityTypeWhere
}

type CalamityTypeWhere {
  id: IDFilter
  name: StringFilter
  order: IntFilter
  image: StringFilter
  userOrganizationId: IDFilter
  batchCreate: BooleanFilter
  description: StringFilter
  defaultM2: FloatFilter
  userOrganization: UserOrganizationWhere
  calamities: CalamityWhere
  updatedAt: IntFilter
  createdAt: IntFilter
  deletedAt: IntFilter
  or: CalamityTypeWhere
  and: CalamityTypeWhere
}

type CompanyFilter {
  search: String
  where: CompanyWhere
}

type CompanyWhere {
  id: IDFilter
  name: StringFilter
  backgroundColor: StringFilter
  textColor: StringFilter
  order: IntFilter
  userOrganizationId: IDFilter
  userOrganization: UserOrganizationWhere
  addresses: AddressWhere
  calamities: CalamityWhere
  companyUsers: CompanyUserWhere
  updatedAt: IntFilter
  createdAt: IntFilter
  deletedAt: IntFilter
  or: CompanyWhere
  and: CompanyWhere
}

type CompanyUserFilter {
  search: String
  where: CompanyUserWhere
}

type CompanyUserWhere {
  id: IDFilter
  userId: IDFilter
  companyId: IDFilter
  company: CompanyWhere
  user: UserWhere
  deletedAt: IntFilter
  updatedAt: IntFilter
  createdAt: IntFilter
  or: CompanyUserWhere
  and: CompanyUserWhere
}

type FailedJobFilter {
  search: String
  where: FailedJobWhere
}

type FailedJobWhere {
  id: IDFilter
  connection: StringFilter
  queue: StringFilter
  payload: StringFilter
  exception: StringFilter
  failedAt: IntFilter
  or: FailedJobWhere
  and: FailedJobWhere
}

type HouseTypeFilter {
  search: String
  where: HouseTypeWhere
}

type HouseTypeWhere {
  id: IDFilter
  name: StringFilter
  order: IntFilter
  image: StringFilter
  userOrganizationId: IDFilter
  description: StringFilter
  userOrganization: UserOrganizationWhere
  addresses: AddressWhere
  calamities: CalamityWhere
  deletedAt: IntFilter
  updatedAt: IntFilter
  createdAt: IntFilter
  or: HouseTypeWhere
  and: HouseTypeWhere
}

type MigrationFilter {
  search: String
  where: MigrationWhere
}

type MigrationWhere {
  id: IDFilter
  migration: StringFilter
  batch: IntFilter
  or: MigrationWhere
  and: MigrationWhere
}

type PersonFilter {
  search: String
  where: PersonWhere
}

type PersonWhere {
  id: IDFilter
  firstName: StringFilter
  lastName: StringFilter
  description: StringFilter
  telephoneNumber: StringFilter
  email: StringFilter
  hasGivenPermission: BooleanFilter
  addressId: IDFilter
  userOrganizationId: IDFilter
  name: StringFilter
  companyName: StringFilter
  address: AddressWhere
  userOrganization: UserOrganizationWhere
  contactPersonAddresses: AddressWhere
  ownerAddresses: AddressWhere
  deletedAt: IntFilter
  createdAt: IntFilter
  updatedAt: IntFilter
  or: PersonWhere
  and: PersonWhere
}

type RoleFilter {
  search: String
  where: RoleWhere
}

type RoleWhere {
  id: IDFilter
  name: StringFilter
  icon: StringFilter
  description: StringFilter
  users: UserWhere
  updatedAt: IntFilter
  deletedAt: IntFilter
  createdAt: IntFilter
  or: RoleWhere
  and: RoleWhere
}

type StatusFilter {
  search: String
  where: StatusWhere
}

type StatusWhere {
  id: IDFilter
  order: IntFilter
  name: StringFilter
  icon: StringFilter
  color: StringFilter
  iconThumb: StringFilter
  calamities: CalamityWhere
  updatedAt: IntFilter
  createdAt: IntFilter
  deletedAt: IntFilter
  or: StatusWhere
  and: StatusWhere
}

type TagFilter {
  search: String
  where: TagWhere
}

type TagWhere {
  id: IDFilter
  name: StringFilter
  image: StringFilter
  userOrganizationId: IDFilter
  userOrganization: UserOrganizationWhere
  calamityTags: CalamityTagWhere
  createdAt: IntFilter
  updatedAt: IntFilter
  deletedAt: IntFilter
  or: TagWhere
  and: TagWhere
}

type UserFilter {
  search: String
  where: UserWhere
}

type UserWhere {
  id: IDFilter
  name: StringFilter
  lastName: StringFilter
  email: StringFilter
  password: StringFilter
  rememberToken: StringFilter
  roleId: IDFilter
  sendNotificationsOnNewCalamity: BooleanFilter
  role: RoleWhere
  calamities: CalamityWhere
  companyUsers: CompanyUserWhere
  userUserOrganizations: UserUserOrganizationWhere
  createdAt: IntFilter
  deletedAt: IntFilter
  updatedAt: IntFilter
  or: UserWhere
  and: UserWhere
}

type UserOrganizationFilter {
  search: String
  where: UserOrganizationWhere
}

type UserOrganizationWhere {
  id: IDFilter
  name: StringFilter
  logo: StringFilter
  street: StringFilter
  city: StringFilter
  houseNumber: StringFilter
  zipCode: StringFilter
  telephoneNumber: StringFilter
  email: StringFilter
  primaryColor: StringFilter
  accentColor: StringFilter
  enablePermissionOnAddress: BooleanFilter
  enableContact: BooleanFilter
  enableOwner: BooleanFilter
  enableContactPerson: BooleanFilter
  newCalamityCompanyScreen: BooleanFilter
  newCalamityPriorityScreen: BooleanFilter
  newCalamityM2Screen: BooleanFilter
  addresses: AddressWhere
  calamities: CalamityWhere
  calamityTypes: CalamityTypeWhere
  companies: CompanyWhere
  houseTypes: HouseTypeWhere
  people: PersonWhere
  tags: TagWhere
  userUserOrganizations: UserUserOrganizationWhere
  createdAt: IntFilter
  updatedAt: IntFilter
  deletedAt: IntFilter
  or: UserOrganizationWhere
  and: UserOrganizationWhere
}

type UserUserOrganizationFilter {
  search: String
  where: UserUserOrganizationWhere
}

type UserUserOrganizationWhere {
  id: IDFilter
  userOrganizationId: IDFilter
  userId: IDFilter
  user: UserWhere
  userOrganization: UserOrganizationWhere
  or: UserUserOrganizationWhere
  and: UserUserOrganizationWhere
}

type Query {
  address(id: ID!): Address!
  addresses(filter: AddressFilter): [Address!]!
  addressStatus(id: ID!): AddressStatus!
  addressStatuses(filter: AddressStatusFilter): [AddressStatus!]!
  calamity(id: ID!): Calamity!
  calamities(filter: CalamityFilter): [Calamity!]!
  calamityAttribute(id: ID!): CalamityAttribute!
  calamityAttributes(filter: CalamityAttributeFilter): [CalamityAttribute!]!
  calamityPicture(id: ID!): CalamityPicture!
  calamityPictures(filter: CalamityPictureFilter): [CalamityPicture!]!
  calamityTag(id: ID!): CalamityTag!
  calamityTags(filter: CalamityTagFilter): [CalamityTag!]!
  calamityType(id: ID!): CalamityType!
  calamityTypes(filter: CalamityTypeFilter): [CalamityType!]!
  company(id: ID!): Company!
  companies(filter: CompanyFilter): [Company!]!
  companyUser(id: ID!): CompanyUser!
  companyUsers(filter: CompanyUserFilter): [CompanyUser!]!
  failedJob(id: ID!): FailedJob!
  failedJobs(filter: FailedJobFilter): [FailedJob!]!
  houseType(id: ID!): HouseType!
  houseTypes(filter: HouseTypeFilter): [HouseType!]!
  migration(id: ID!): Migration!
  migrations(filter: MigrationFilter): [Migration!]!
  person(id: ID!): Person!
  people(filter: PersonFilter): [Person!]!
  role(id: ID!): Role!
  roles(filter: RoleFilter): [Role!]!
  status(id: ID!): Status!
  statuses(filter: StatusFilter): [Status!]!
  tag(id: ID!): Tag!
  tags(filter: TagFilter): [Tag!]!
  user(id: ID!): User!
  users(filter: UserFilter): [User!]!
  userOrganization(id: ID!): UserOrganization!
  userOrganizations(filter: UserOrganizationFilter): [UserOrganization!]!
  userUserOrganization(id: ID!): UserUserOrganization!
  userUserOrganizations(
    filter: UserUserOrganizationFilter
  ): [UserUserOrganization!]!
}

input AddressInput {
  id: ID!
  street: String
  houseNumber: String
  zipAddress: String
  city: String
  longitude: Float!
  latitude: Float!
  ownerId: ID
  contactPersonId: ID
  userOrganizationId: ID!
  companyId: ID
  houseTypeId: ID
  description: String
  name: String
  permission: Boolean
  addressStatusId: ID!
  createdAt: Int
  deletedAt: Int
  updatedAt: Int
}

type AddressPayload {
  address: Address!
}

type AddressDeletePayload {
  id: ID!
}

type AddressesDeletePayload {
  ids: [ID!]!
}

type AddressesUpdatePayload {
  ok: Boolean!
}

input AddressStatusInput {
  id: ID!
  name: String
  icon: String
  description: String
  order: Int!
  color: String!
  createdAt: Int
  updatedAt: Int
  deletedAt: Int
}

type AddressStatusPayload {
  addressStatus: AddressStatus!
}

type AddressStatusDeletePayload {
  id: ID!
}

type AddressStatusesDeletePayload {
  ids: [ID!]!
}

type AddressStatusesUpdatePayload {
  ok: Boolean!
}

input CalamityInput {
  id: ID!
  parentId: ID
  description: String
  priority: Int!
  totalBuildings: String
  street: String
  houseNumbers: String
  zipAddress: String
  longitude: String
  latitude: String
  statusId: ID
  companyId: ID
  userId: ID
  houseTypeId: ID
  calamityTypeId: ID
  city: String
  userOrganizationId: ID!
  addressId: ID
  underground: String
  part: String
  color: String
  whatColor: String
  resolvedM2: Float
  resolvedDate: Int
  resolvedCompanyDescription: String
  resolvedUserDescription: String
  m2: Float
  amount: Int
  createdAt: Int
  updatedAt: Int
  deletedAt: Int
}

type CalamityPayload {
  calamity: Calamity!
}

type CalamityDeletePayload {
  id: ID!
}

type CalamitiesDeletePayload {
  ids: [ID!]!
}

type CalamitiesUpdatePayload {
  ok: Boolean!
}

input CalamityAttributeInput {
  id: ID!
  key: String!
  type: String!
  value: String
  calamityId: ID!
  deletedAt: Int
  createdAt: Int
  updatedAt: Int
}

type CalamityAttributePayload {
  calamityAttribute: CalamityAttribute!
}

type CalamityAttributeDeletePayload {
  id: ID!
}

type CalamityAttributesDeletePayload {
  ids: [ID!]!
}

type CalamityAttributesUpdatePayload {
  ok: Boolean!
}

input CalamityPictureInput {
  id: ID!
  name: String
  url: String!
  calamityId: ID!
  deletedAt: Int
  createdAt: Int
  updatedAt: Int
}

type CalamityPicturePayload {
  calamityPicture: CalamityPicture!
}

type CalamityPictureDeletePayload {
  id: ID!
}

type CalamityPicturesDeletePayload {
  ids: [ID!]!
}

type CalamityPicturesUpdatePayload {
  ok: Boolean!
}

input CalamityTagInput {
  id: ID!
  calamityId: ID!
  tagId: ID!
}

type CalamityTagPayload {
  calamityTag: CalamityTag!
}

type CalamityTagDeletePayload {
  id: ID!
}

type CalamityTagsDeletePayload {
  ids: [ID!]!
}

type CalamityTagsUpdatePayload {
  ok: Boolean!
}

input CalamityTypeInput {
  id: ID!
  name: String!
  order: Int!
  image: String
  userOrganizationId: ID!
  batchCreate: Boolean
  description: String
  defaultM2: Float!
  updatedAt: Int
  createdAt: Int
  deletedAt: Int
}

type CalamityTypePayload {
  calamityType: CalamityType!
}

type CalamityTypeDeletePayload {
  id: ID!
}

type CalamityTypesDeletePayload {
  ids: [ID!]!
}

type CalamityTypesUpdatePayload {
  ok: Boolean!
}

input CompanyInput {
  id: ID!
  name: String!
  backgroundColor: String!
  textColor: String!
  order: Int
  userOrganizationId: ID!
  updatedAt: Int
  createdAt: Int
  deletedAt: Int
}

type CompanyPayload {
  company: Company!
}

type CompanyDeletePayload {
  id: ID!
}

type CompaniesDeletePayload {
  ids: [ID!]!
}

type CompaniesUpdatePayload {
  ok: Boolean!
}

input CompanyUserInput {
  id: ID!
  userId: ID!
  companyId: ID!
  deletedAt: Int
  updatedAt: Int
  createdAt: Int
}

type CompanyUserPayload {
  companyUser: CompanyUser!
}

type CompanyUserDeletePayload {
  id: ID!
}

type CompanyUsersDeletePayload {
  ids: [ID!]!
}

type CompanyUsersUpdatePayload {
  ok: Boolean!
}

input FailedJobInput {
  id: ID!
  connection: String!
  queue: String!
  payload: String!
  exception: String!
  failedAt: Int!
}

type FailedJobPayload {
  failedJob: FailedJob!
}

type FailedJobDeletePayload {
  id: ID!
}

type FailedJobsDeletePayload {
  ids: [ID!]!
}

type FailedJobsUpdatePayload {
  ok: Boolean!
}

input HouseTypeInput {
  id: ID!
  name: String!
  order: Int!
  image: String
  userOrganizationId: ID!
  description: String
  deletedAt: Int
  updatedAt: Int
  createdAt: Int
}

type HouseTypePayload {
  houseType: HouseType!
}

type HouseTypeDeletePayload {
  id: ID!
}

type HouseTypesDeletePayload {
  ids: [ID!]!
}

type HouseTypesUpdatePayload {
  ok: Boolean!
}

input MigrationInput {
  id: ID!
  migration: String!
  batch: Int!
}

type MigrationPayload {
  migration: Migration!
}

type MigrationDeletePayload {
  id: ID!
}

type MigrationsDeletePayload {
  ids: [ID!]!
}

type MigrationsUpdatePayload {
  ok: Boolean!
}

input PersonInput {
  id: ID!
  firstName: String
  lastName: String
  description: String
  telephoneNumber: String
  email: String
  hasGivenPermission: Boolean
  addressId: ID
  userOrganizationId: ID!
  name: String
  companyName: String
  deletedAt: Int
  createdAt: Int
  updatedAt: Int
}

type PersonPayload {
  person: Person!
}

type PersonDeletePayload {
  id: ID!
}

type PeopleDeletePayload {
  ids: [ID!]!
}

type PeopleUpdatePayload {
  ok: Boolean!
}

input RoleInput {
  id: ID!
  name: String!
  icon: String
  description: String
  updatedAt: Int
  deletedAt: Int
  createdAt: Int
}

type RolePayload {
  role: Role!
}

type RoleDeletePayload {
  id: ID!
}

type RolesDeletePayload {
  ids: [ID!]!
}

type RolesUpdatePayload {
  ok: Boolean!
}

input StatusInput {
  id: ID!
  order: Int!
  name: String!
  icon: String
  color: String
  iconThumb: String
  updatedAt: Int
  createdAt: Int
  deletedAt: Int
}

type StatusPayload {
  status: Status!
}

type StatusDeletePayload {
  id: ID!
}

type StatusesDeletePayload {
  ids: [ID!]!
}

type StatusesUpdatePayload {
  ok: Boolean!
}

input TagInput {
  id: ID!
  name: String!
  image: String
  userOrganizationId: ID!
  createdAt: Int
  updatedAt: Int
  deletedAt: Int
}

type TagPayload {
  tag: Tag!
}

type TagDeletePayload {
  id: ID!
}

type TagsDeletePayload {
  ids: [ID!]!
}

type TagsUpdatePayload {
  ok: Boolean!
}

input UserInput {
  id: ID!
  name: String!
  lastName: String!
  email: String!
  password: String!
  rememberToken: String
  roleId: ID!
  sendNotificationsOnNewCalamity: Boolean!
  createdAt: Int
  deletedAt: Int
  updatedAt: Int
}

type UserPayload {
  user: User!
}

type UserDeletePayload {
  id: ID!
}

type UsersDeletePayload {
  ids: [ID!]!
}

type UsersUpdatePayload {
  ok: Boolean!
}

input UserOrganizationInput {
  id: ID!
  name: String
  logo: String
  street: String
  city: String
  houseNumber: String
  zipCode: String
  telephoneNumber: String
  email: String
  primaryColor: String
  accentColor: String
  enablePermissionOnAddress: Boolean!
  enableContact: Boolean!
  enableOwner: Boolean!
  enableContactPerson: Boolean!
  newCalamityCompanyScreen: Boolean!
  newCalamityPriorityScreen: Boolean!
  newCalamityM2Screen: Boolean!
  createdAt: Int
  updatedAt: Int
  deletedAt: Int
}

type UserOrganizationPayload {
  userOrganization: UserOrganization!
}

type UserOrganizationDeletePayload {
  id: ID!
}

type UserOrganizationsDeletePayload {
  ids: [ID!]!
}

type UserOrganizationsUpdatePayload {
  ok: Boolean!
}

input UserUserOrganizationInput {
  id: ID!
  userOrganizationId: ID!
  userId: ID!
}

type UserUserOrganizationPayload {
  userUserOrganization: UserUserOrganization!
}

type UserUserOrganizationDeletePayload {
  id: ID!
}

type UserUserOrganizationsDeletePayload {
  ids: [ID!]!
}

type UserUserOrganizationsUpdatePayload {
  ok: Boolean!
}

type Mutation {
  createAddress(input: AddressInput!): AddressPayload!
  createAddresses(input: AddressesInput!): AddressesPayload!
  updateAddress(input: AddressInput!): AddressPayload!
  updateAddresses(
    filter: AddressFilter
    input: AddressesInput!
  ): AddressesUpdatePayload!
  deleteAddress(id: ID!): AddressDeletePayload!
  deleteAddresses(filter: AddressFilter): AddressesDeletePayload!
  createAddressStatus(input: AddressStatusInput!): AddressStatusPayload!
  createAddressStatuses(input: AddressStatusesInput!): AddressStatusesPayload!
  updateAddressStatus(input: AddressStatusInput!): AddressStatusPayload!
  updateAddressStatuses(
    filter: AddressStatusFilter
    input: AddressStatusesInput!
  ): AddressStatusesUpdatePayload!
  deleteAddressStatus(id: ID!): AddressStatusDeletePayload!
  deleteAddressStatuses(
    filter: AddressStatusFilter
  ): AddressStatusesDeletePayload!
  createCalamity(input: CalamityInput!): CalamityPayload!
  createCalamities(input: CalamitiesInput!): CalamitiesPayload!
  updateCalamity(input: CalamityInput!): CalamityPayload!
  updateCalamities(
    filter: CalamityFilter
    input: CalamitiesInput!
  ): CalamitiesUpdatePayload!
  deleteCalamity(id: ID!): CalamityDeletePayload!
  deleteCalamities(filter: CalamityFilter): CalamitiesDeletePayload!
  createCalamityAttribute(
    input: CalamityAttributeInput!
  ): CalamityAttributePayload!
  createCalamityAttributes(
    input: CalamityAttributesInput!
  ): CalamityAttributesPayload!
  updateCalamityAttribute(
    input: CalamityAttributeInput!
  ): CalamityAttributePayload!
  updateCalamityAttributes(
    filter: CalamityAttributeFilter
    input: CalamityAttributesInput!
  ): CalamityAttributesUpdatePayload!
  deleteCalamityAttribute(id: ID!): CalamityAttributeDeletePayload!
  deleteCalamityAttributes(
    filter: CalamityAttributeFilter
  ): CalamityAttributesDeletePayload!
  createCalamityPicture(input: CalamityPictureInput!): CalamityPicturePayload!
  createCalamityPictures(
    input: CalamityPicturesInput!
  ): CalamityPicturesPayload!
  updateCalamityPicture(input: CalamityPictureInput!): CalamityPicturePayload!
  updateCalamityPictures(
    filter: CalamityPictureFilter
    input: CalamityPicturesInput!
  ): CalamityPicturesUpdatePayload!
  deleteCalamityPicture(id: ID!): CalamityPictureDeletePayload!
  deleteCalamityPictures(
    filter: CalamityPictureFilter
  ): CalamityPicturesDeletePayload!
  createCalamityTag(input: CalamityTagInput!): CalamityTagPayload!
  createCalamityTags(input: CalamityTagsInput!): CalamityTagsPayload!
  updateCalamityTag(input: CalamityTagInput!): CalamityTagPayload!
  updateCalamityTags(
    filter: CalamityTagFilter
    input: CalamityTagsInput!
  ): CalamityTagsUpdatePayload!
  deleteCalamityTag(id: ID!): CalamityTagDeletePayload!
  deleteCalamityTags(filter: CalamityTagFilter): CalamityTagsDeletePayload!
  createCalamityType(input: CalamityTypeInput!): CalamityTypePayload!
  createCalamityTypes(input: CalamityTypesInput!): CalamityTypesPayload!
  updateCalamityType(input: CalamityTypeInput!): CalamityTypePayload!
  updateCalamityTypes(
    filter: CalamityTypeFilter
    input: CalamityTypesInput!
  ): CalamityTypesUpdatePayload!
  deleteCalamityType(id: ID!): CalamityTypeDeletePayload!
  deleteCalamityTypes(filter: CalamityTypeFilter): CalamityTypesDeletePayload!
  createCompany(input: CompanyInput!): CompanyPayload!
  createCompanies(input: CompaniesInput!): CompaniesPayload!
  updateCompany(input: CompanyInput!): CompanyPayload!
  updateCompanies(
    filter: CompanyFilter
    input: CompaniesInput!
  ): CompaniesUpdatePayload!
  deleteCompany(id: ID!): CompanyDeletePayload!
  deleteCompanies(filter: CompanyFilter): CompaniesDeletePayload!
  createCompanyUser(input: CompanyUserInput!): CompanyUserPayload!
  createCompanyUsers(input: CompanyUsersInput!): CompanyUsersPayload!
  updateCompanyUser(input: CompanyUserInput!): CompanyUserPayload!
  updateCompanyUsers(
    filter: CompanyUserFilter
    input: CompanyUsersInput!
  ): CompanyUsersUpdatePayload!
  deleteCompanyUser(id: ID!): CompanyUserDeletePayload!
  deleteCompanyUsers(filter: CompanyUserFilter): CompanyUsersDeletePayload!
  createFailedJob(input: FailedJobInput!): FailedJobPayload!
  createFailedJobs(input: FailedJobsInput!): FailedJobsPayload!
  updateFailedJob(input: FailedJobInput!): FailedJobPayload!
  updateFailedJobs(
    filter: FailedJobFilter
    input: FailedJobsInput!
  ): FailedJobsUpdatePayload!
  deleteFailedJob(id: ID!): FailedJobDeletePayload!
  deleteFailedJobs(filter: FailedJobFilter): FailedJobsDeletePayload!
  createHouseType(input: HouseTypeInput!): HouseTypePayload!
  createHouseTypes(input: HouseTypesInput!): HouseTypesPayload!
  updateHouseType(input: HouseTypeInput!): HouseTypePayload!
  updateHouseTypes(
    filter: HouseTypeFilter
    input: HouseTypesInput!
  ): HouseTypesUpdatePayload!
  deleteHouseType(id: ID!): HouseTypeDeletePayload!
  deleteHouseTypes(filter: HouseTypeFilter): HouseTypesDeletePayload!
  createMigration(input: MigrationInput!): MigrationPayload!
  createMigrations(input: MigrationsInput!): MigrationsPayload!
  updateMigration(input: MigrationInput!): MigrationPayload!
  updateMigrations(
    filter: MigrationFilter
    input: MigrationsInput!
  ): MigrationsUpdatePayload!
  deleteMigration(id: ID!): MigrationDeletePayload!
  deleteMigrations(filter: MigrationFilter): MigrationsDeletePayload!
  createPerson(input: PersonInput!): PersonPayload!
  createPeople(input: PeopleInput!): PeoplePayload!
  updatePerson(input: PersonInput!): PersonPayload!
  updatePeople(filter: PersonFilter, input: PeopleInput!): PeopleUpdatePayload!
  deletePerson(id: ID!): PersonDeletePayload!
  deletePeople(filter: PersonFilter): PeopleDeletePayload!
  createRole(input: RoleInput!): RolePayload!
  createRoles(input: RolesInput!): RolesPayload!
  updateRole(input: RoleInput!): RolePayload!
  updateRoles(filter: RoleFilter, input: RolesInput!): RolesUpdatePayload!
  deleteRole(id: ID!): RoleDeletePayload!
  deleteRoles(filter: RoleFilter): RolesDeletePayload!
  createStatus(input: StatusInput!): StatusPayload!
  createStatuses(input: StatusesInput!): StatusesPayload!
  updateStatus(input: StatusInput!): StatusPayload!
  updateStatuses(
    filter: StatusFilter
    input: StatusesInput!
  ): StatusesUpdatePayload!
  deleteStatus(id: ID!): StatusDeletePayload!
  deleteStatuses(filter: StatusFilter): StatusesDeletePayload!
  createTag(input: TagInput!): TagPayload!
  createTags(input: TagsInput!): TagsPayload!
  updateTag(input: TagInput!): TagPayload!
  updateTags(filter: TagFilter, input: TagsInput!): TagsUpdatePayload!
  deleteTag(id: ID!): TagDeletePayload!
  deleteTags(filter: TagFilter): TagsDeletePayload!
  createUser(input: UserInput!): UserPayload!
  createUsers(input: UsersInput!): UsersPayload!
  updateUser(input: UserInput!): UserPayload!
  updateUsers(filter: UserFilter, input: UsersInput!): UsersUpdatePayload!
  deleteUser(id: ID!): UserDeletePayload!
  deleteUsers(filter: UserFilter): UsersDeletePayload!
  createUserOrganization(
    input: UserOrganizationInput!
  ): UserOrganizationPayload!
  createUserOrganizations(
    input: UserOrganizationsInput!
  ): UserOrganizationsPayload!
  updateUserOrganization(
    input: UserOrganizationInput!
  ): UserOrganizationPayload!
  updateUserOrganizations(
    filter: UserOrganizationFilter
    input: UserOrganizationsInput!
  ): UserOrganizationsUpdatePayload!
  deleteUserOrganization(id: ID!): UserOrganizationDeletePayload!
  deleteUserOrganizations(
    filter: UserOrganizationFilter
  ): UserOrganizationsDeletePayload!
  createUserUserOrganization(
    input: UserUserOrganizationInput!
  ): UserUserOrganizationPayload!
  createUserUserOrganizations(
    input: UserUserOrganizationsInput!
  ): UserUserOrganizationsPayload!
  updateUserUserOrganization(
    input: UserUserOrganizationInput!
  ): UserUserOrganizationPayload!
  updateUserUserOrganizations(
    filter: UserUserOrganizationFilter
    input: UserUserOrganizationsInput!
  ): UserUserOrganizationsUpdatePayload!
  deleteUserUserOrganization(id: ID!): UserUserOrganizationDeletePayload!
  deleteUserUserOrganizations(
    filter: UserUserOrganizationFilter
  ): UserUserOrganizationsDeletePayload!
}

```

## Filtering lists (WIP)

This program generates type safe filters you can use in your frontend

### Search

```graphql
query(filter: {
    search: 'jan'
})
```

### This or that

```graphql
query(filter: {
    where: {
        name:{
            equalTo: 'Jan'
        }
        or: {
            name: {
                equalTo: 'Jannes',
            }
        }
    }
})

where: {
    id: {
        equalTo: 1
    }
    name: {
        startsWith: 'J',
    }
    where:{

    }
}
```

### (() or ())

````graphql
query(filter: {
    where: {
        or:{
            id: {
                equalTo: 1
            }
            name: {
                startsWith: 'J',
            }
            or:{
                id: {
                    equalTo: 2
                }
                name: {
                    startsWith: 'R',
                }
            }
        }
    }})

### Filter
```graphql
query(filter: {
    where: {
        id:{
            in: [1, 2, 3, 4]
            equalTo: 1
        },
        name: {
            like: "joe"
        },
        organizationId: {
            equalTo: 1,
        }
        or: {
            id:{
                in: []
                equalTo: 1
            },
             name: {
                like: "joe"
            },
        }
        and: {

        }
    }

})
````

## How to detect many to many (notes to myself)

```golang
User {
    UserOrganizations []UserOrganization
    Posts []Post
}

Post {
    User User
    Message string
}

UserOrganization {
    User User
    Organization Organization
}

Organization {
    Users []User
}
```

### Possible detection

if table has only 2 relationship and table contains both tables
