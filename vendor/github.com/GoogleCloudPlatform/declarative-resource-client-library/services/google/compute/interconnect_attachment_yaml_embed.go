// Copyright 2023 Google LLC. All Rights Reserved.
// 
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// 
//     http://www.apache.org/licenses/LICENSE-2.0
// 
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
// GENERATED BY gen_go_data.go
// gen_go_data -package compute -var YAML_interconnect_attachment blaze-out/k8-fastbuild/genfiles/cloud/graphite/mmv2/services/google/compute/interconnect_attachment.yaml

package compute

// blaze-out/k8-fastbuild/genfiles/cloud/graphite/mmv2/services/google/compute/interconnect_attachment.yaml
var YAML_interconnect_attachment = []byte("info:\n  title: Compute/InterconnectAttachment\n  description: The Compute InterconnectAttachment resource\n  x-dcl-struct-name: InterconnectAttachment\n  x-dcl-has-iam: false\npaths:\n  get:\n    description: The function used to get information about a InterconnectAttachment\n    parameters:\n    - name: interconnectAttachment\n      required: true\n      description: A full instance of a InterconnectAttachment\n  apply:\n    description: The function used to apply information about a InterconnectAttachment\n    parameters:\n    - name: interconnectAttachment\n      required: true\n      description: A full instance of a InterconnectAttachment\n  delete:\n    description: The function used to delete a InterconnectAttachment\n    parameters:\n    - name: interconnectAttachment\n      required: true\n      description: A full instance of a InterconnectAttachment\n  deleteAll:\n    description: The function used to delete all InterconnectAttachment\n    parameters:\n    - name: project\n      required: true\n      schema:\n        type: string\n    - name: region\n      required: true\n      schema:\n        type: string\n  list:\n    description: The function used to list information about many InterconnectAttachment\n    parameters:\n    - name: project\n      required: true\n      schema:\n        type: string\n    - name: region\n      required: true\n      schema:\n        type: string\ncomponents:\n  schemas:\n    InterconnectAttachment:\n      title: InterconnectAttachment\n      x-dcl-id: projects/{{project}}/regions/{{region}}/interconnectAttachments/{{name}}\n      x-dcl-locations:\n      - region\n      x-dcl-parent-container: project\n      x-dcl-has-create: true\n      x-dcl-has-iam: false\n      x-dcl-read-timeout: 0\n      x-dcl-apply-timeout: 0\n      x-dcl-delete-timeout: 0\n      type: object\n      required:\n      - name\n      - region\n      - project\n      properties:\n        adminEnabled:\n          type: boolean\n          x-dcl-go-name: AdminEnabled\n          description: Determines whether this Attachment will carry packets. Not\n            present for PARTNER_PROVIDER.\n        bandwidth:\n          type: string\n          x-dcl-go-name: Bandwidth\n          x-dcl-go-type: InterconnectAttachmentBandwidthEnum\n          description: 'Provisioned bandwidth capacity for the interconnect attachment.\n            For attachments of type DEDICATED, the user can set the bandwidth. For\n            attachments of type PARTNER, the Google Partner that is operating the\n            interconnect must set the bandwidth. Output only for PARTNER type, mutable\n            for PARTNER_PROVIDER and DEDICATED, and can take one of the following\n            values: - BPS_50M: 50 Mbit/s - BPS_100M: 100 Mbit/s - BPS_200M: 200 Mbit/s\n            - BPS_300M: 300 Mbit/s - BPS_400M: 400 Mbit/s - BPS_500M: 500 Mbit/s -\n            BPS_1G: 1 Gbit/s - BPS_2G: 2 Gbit/s - BPS_5G: 5 Gbit/s - BPS_10G: 10 Gbit/s\n            - BPS_20G: 20 Gbit/s - BPS_50G: 50 Gbit/s'\n          enum:\n          - BPS_50M\n          - BPS_100M\n          - BPS_200M\n          - BPS_300M\n          - BPS_400M\n          - BPS_500M\n          - BPS_1G\n          - BPS_2G\n          - BPS_5G\n          - BPS_10G\n          - BPS_20G\n          - BPS_50G\n        candidateSubnets:\n          type: array\n          x-dcl-go-name: CandidateSubnets\n          description: Up to 16 candidate prefixes that can be used to restrict the\n            allocation of cloudRouterIpAddress and customerRouterIpAddress for this\n            attachment. All prefixes must be within link-local address space (169.254.0.0/16)\n            and must be /29 or shorter (/28, /27, etc). Google will attempt to select\n            an unused /29 from the supplied candidate prefix(es). The request will\n            fail if all possible /29s are in use on Google's edge. If not supplied,\n            Google will randomly select an unused /29 from all of link-local space.\n          x-dcl-send-empty: true\n          x-dcl-list-type: list\n          items:\n            type: string\n            x-dcl-go-type: string\n        cloudRouterIPAddress:\n          type: string\n          x-dcl-go-name: CloudRouterIPAddress\n          readOnly: true\n          description: IPv4 address + prefix length to be configured on Cloud Router\n            Interface for this interconnect attachment.\n          x-kubernetes-immutable: true\n        customerRouterIPAddress:\n          type: string\n          x-dcl-go-name: CustomerRouterIPAddress\n          readOnly: true\n          description: IPv4 address + prefix length to be configured on the customer\n            router subinterface for this interconnect attachment.\n          x-kubernetes-immutable: true\n        dataplaneVersion:\n          type: integer\n          format: int64\n          x-dcl-go-name: DataplaneVersion\n          description: Dataplane version for this InterconnectAttachment. This field\n            is only present for Dataplane version 2 and higher. Absence of this field\n            in the API output indicates that the Dataplane is version 1.\n        description:\n          type: string\n          x-dcl-go-name: Description\n          description: An optional description of this resource.\n        edgeAvailabilityDomain:\n          type: string\n          x-dcl-go-name: EdgeAvailabilityDomain\n          x-dcl-go-type: InterconnectAttachmentEdgeAvailabilityDomainEnum\n          description: 'Desired availability domain for the attachment. Only available\n            for type PARTNER, at creation time, and can take one of the following\n            values: - AVAILABILITY_DOMAIN_ANY - AVAILABILITY_DOMAIN_1 - AVAILABILITY_DOMAIN_2\n            For improved reliability, customers should configure a pair of attachments,\n            one per availability domain. The selected availability domain will be\n            provided to the Partner via the pairing key, so that the provisioned circuit\n            will lie in the specified domain. If not specified, the value will default\n            to AVAILABILITY_DOMAIN_ANY.'\n          enum:\n          - AVAILABILITY_DOMAIN_ANY\n          - AVAILABILITY_DOMAIN_1\n          - AVAILABILITY_DOMAIN_2\n        encryption:\n          type: string\n          x-dcl-go-name: Encryption\n          x-dcl-go-type: InterconnectAttachmentEncryptionEnum\n          description: 'Indicates the user-supplied encryption option of this VLAN\n            attachment (interconnectAttachment). Can only be specified at attachment\n            creation for PARTNER or DEDICATED attachments. Possible values are: -\n            `NONE` - This is the default value, which means that the VLAN attachment\n            carries unencrypted traffic. VMs are able to send traffic to, or receive\n            traffic from, such a VLAN attachment. - `IPSEC` - The VLAN attachment\n            carries only encrypted traffic that is encrypted by an IPsec device, such\n            as an HA VPN gateway or third-party IPsec VPN. VMs cannot directly send\n            traffic to, or receive traffic from, such a VLAN attachment. To use _IPsec-encrypted\n            Cloud Interconnect_, the VLAN attachment must be created with this option.\n            Not currently available publicly.'\n          enum:\n          - NONE\n          - IPSEC\n        id:\n          type: integer\n          format: int64\n          x-dcl-go-name: Id\n          readOnly: true\n          description: The unique identifier for the resource. This identifier is\n            defined by the server.\n          x-kubernetes-immutable: true\n        interconnect:\n          type: string\n          x-dcl-go-name: Interconnect\n          description: URL of the underlying Interconnect object that this attachment's\n            traffic will traverse through.\n        ipsecInternalAddresses:\n          type: array\n          x-dcl-go-name: IpsecInternalAddresses\n          description: A list of URLs of addresses that have been reserved for the\n            VLAN attachment. Used only for the VLAN attachment that has the encryption\n            option as IPSEC. The addresses must be regional internal IP address ranges.\n            When creating an HA VPN gateway over the VLAN attachment, if the attachment\n            is configured to use a regional internal IP address, then the VPN gateway's\n            IP address is allocated from the IP address range specified here. For\n            example, if the HA VPN gateway's interface 0 is paired to this VLAN attachment,\n            then a regional internal IP address for the VPN gateway interface 0 will\n            be allocated from the IP address specified for this VLAN attachment. If\n            this field is not specified when creating the VLAN attachment, then later\n            on when creating an HA VPN gateway on this VLAN attachment, the HA VPN\n            gateway's IP address is allocated from the regional external IP address\n            pool. Not currently available publicly.\n          x-dcl-send-empty: true\n          x-dcl-list-type: list\n          items:\n            type: string\n            x-dcl-go-type: string\n        mtu:\n          type: integer\n          format: int64\n          x-dcl-go-name: Mtu\n          description: Maximum Transmission Unit (MTU), in bytes, of packets passing\n            through this interconnect attachment. Only 1440 and 1500 are allowed.\n            If not specified, the value will default to 1440.\n        name:\n          type: string\n          x-dcl-go-name: Name\n          description: Name of the resource. Provided by the client when the resource\n            is created. The name must be 1-63 characters long, and comply with [RFC1035](https://www.ietf.org/rfc/rfc1035.txt).\n            Specifically, the name must be 1-63 characters long and match the regular\n            expression `[a-z]([-a-z0-9]*[a-z0-9])?` which means the first character\n            must be a lowercase letter, and all following characters must be a dash,\n            lowercase letter, or digit, except the last character, which cannot be\n            a dash.\n        operationalStatus:\n          type: string\n          x-dcl-go-name: OperationalStatus\n          x-dcl-go-type: InterconnectAttachmentOperationalStatusEnum\n          readOnly: true\n          description: 'The current status of whether or not this interconnect attachment\n            is functional, which can take one of the following values: - OS_ACTIVE:\n            The attachment has been turned up and is ready to use. - OS_UNPROVISIONED:\n            The attachment is not ready to use yet, because turnup is not complete.'\n          x-kubernetes-immutable: true\n          enum:\n          - OS_ACTIVE\n          - OS_UNPROVISIONED\n        pairingKey:\n          type: string\n          x-dcl-go-name: PairingKey\n          description: The opaque identifier of an PARTNER attachment used to initiate\n            provisioning with a selected partner. Of the form \"XXXXX/region/domain\"\n        partnerAsn:\n          type: integer\n          format: int64\n          x-dcl-go-name: PartnerAsn\n          description: Optional BGP ASN for the router supplied by a Layer 3 Partner\n            if they configured BGP on behalf of the customer. Output only for PARTNER\n            type, input only for PARTNER_PROVIDER, not available for DEDICATED.\n        partnerMetadata:\n          type: object\n          x-dcl-go-name: PartnerMetadata\n          x-dcl-go-type: InterconnectAttachmentPartnerMetadata\n          description: Informational metadata about Partner attachments from Partners\n            to display to customers. Output only for for PARTNER type, mutable for\n            PARTNER_PROVIDER, not available for DEDICATED.\n          properties:\n            interconnectName:\n              type: string\n              x-dcl-go-name: InterconnectName\n              description: Plain text name of the Interconnect this attachment is\n                connected to, as displayed in the Partner's portal. For instance \"Chicago\n                1\". This value may be validated to match approved Partner values.\n            partnerName:\n              type: string\n              x-dcl-go-name: PartnerName\n              description: Plain text name of the Partner providing this attachment.\n                This value may be validated to match approved Partner values.\n            portalUrl:\n              type: string\n              x-dcl-go-name: PortalUrl\n              description: URL of the Partner's portal for this Attachment. Partners\n                may customise this to be a deep link to the specific resource on the\n                Partner portal. This value may be validated to match approved Partner\n                values.\n        privateInterconnectInfo:\n          type: object\n          x-dcl-go-name: PrivateInterconnectInfo\n          x-dcl-go-type: InterconnectAttachmentPrivateInterconnectInfo\n          readOnly: true\n          description: Information specific to an InterconnectAttachment. This property\n            is populated if the interconnect that this is attached to is of type DEDICATED.\n          x-kubernetes-immutable: true\n          properties:\n            tag8021q:\n              type: integer\n              format: int64\n              x-dcl-go-name: Tag8021q\n              readOnly: true\n              description: 802.1q encapsulation tag to be used for traffic between\n                Google and the customer, going to and from this network and region.\n              x-kubernetes-immutable: true\n        project:\n          type: string\n          x-dcl-go-name: Project\n          description: The project for the resource\n          x-kubernetes-immutable: true\n          x-dcl-references:\n          - resource: Cloudresourcemanager/Project\n            field: name\n            parent: true\n        region:\n          type: string\n          x-dcl-go-name: Region\n          description: URL of the region where the regional interconnect attachment\n            resides. You must specify this field as part of the HTTP request URL.\n            It is not settable as a field in the request body.\n          x-kubernetes-immutable: true\n        router:\n          type: string\n          x-dcl-go-name: Router\n          description: URL of the Cloud Router to be used for dynamic routing. This\n            router must be in the same region as this InterconnectAttachment. The\n            InterconnectAttachment will automatically connect the Interconnect to\n            the network & region within which the Cloud Router is configured.\n        satisfiesPzs:\n          type: boolean\n          x-dcl-go-name: SatisfiesPzs\n          readOnly: true\n          description: Set to true if the resource satisfies the zone separation organization\n            policy constraints and false otherwise. Defaults to false if the field\n            is not present.\n          x-kubernetes-immutable: true\n        selfLink:\n          type: string\n          x-dcl-go-name: SelfLink\n          readOnly: true\n          description: Server-defined URL for the resource.\n          x-kubernetes-immutable: true\n        state:\n          type: string\n          x-dcl-go-name: State\n          x-dcl-go-type: InterconnectAttachmentStateEnum\n          readOnly: true\n          description: 'The current state of this attachment''s functionality. Enum\n            values ACTIVE and UNPROVISIONED are shared by DEDICATED/PRIVATE, PARTNER,\n            and PARTNER_PROVIDER interconnect attachments, while enum values PENDING_PARTNER,\n            PARTNER_REQUEST_RECEIVED, and PENDING_CUSTOMER are used for only PARTNER\n            and PARTNER_PROVIDER interconnect attachments. This state can take one\n            of the following values: - ACTIVE: The attachment has been turned up and\n            is ready to use. - UNPROVISIONED: The attachment is not ready to use yet,\n            because turnup is not complete. - PENDING_PARTNER: A newly-created PARTNER\n            attachment that has not yet been configured on the Partner side. - PARTNER_REQUEST_RECEIVED:\n            A PARTNER attachment is in the process of provisioning after a PARTNER_PROVIDER\n            attachment was created that references it. - PENDING_CUSTOMER: A PARTNER\n            or PARTNER_PROVIDER attachment that is waiting for a customer to activate\n            it. - DEFUNCT: The attachment was deleted externally and is no longer\n            functional. This could be because the associated Interconnect was removed,\n            or because the other side of a Partner attachment was deleted. Possible\n            values: DEPRECATED, OBSOLETE, DELETED, ACTIVE'\n          x-kubernetes-immutable: true\n          enum:\n          - DEPRECATED\n          - OBSOLETE\n          - DELETED\n          - ACTIVE\n        type:\n          type: string\n          x-dcl-go-name: Type\n          x-dcl-go-type: InterconnectAttachmentTypeEnum\n          description: 'The type of interconnect attachment this is, which can take\n            one of the following values: - DEDICATED: an attachment to a Dedicated\n            Interconnect. - PARTNER: an attachment to a Partner Interconnect, created\n            by the customer. - PARTNER_PROVIDER: an attachment to a Partner Interconnect,\n            created by the partner. Possible values: PATH, OTHER, PARAMETER'\n          enum:\n          - PATH\n          - OTHER\n          - PARAMETER\n        vlanTag8021q:\n          type: integer\n          format: int64\n          x-dcl-go-name: VlanTag8021q\n          description: The IEEE 802.1Q VLAN tag for this attachment, in the range\n            2-4094. Only specified at creation time.\n")

// 17490 bytes
// MD5: 03bc17de2050c53301a9c1165d204c09
