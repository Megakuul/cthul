/**
 * Cthul System
 *
 * Copyright (C) 2024 Linus Ilian Moser <linus.moser@megakuul.ch>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program. If not, see <https://www.gnu.org/licenses/>.
 */

package libvirt

import (
	"context"
	"encoding/xml"
	"fmt"

	"github.com/digitalocean/go-libvirt"
	"cthul.io/cthul/pkg/adapter/domain/libvirt/hotplug"
	"cthul.io/cthul/pkg/adapter/domain/structure"
)

// List fetches the uuid & name of all domains located on this node.
func (l *Adapter) List(ctx context.Context) (map[string]string, error) {
	err := l.initClient()
	if err!=nil {
		return nil, err
	}

	domains, _, err := l.client.ConnectListAllDomains(-1, 0)
	if err!=nil {
		return nil, err
	}

	domainUUIDs := map[string]string{}
	for _, domain := range domains {
		uuid, err := l.serializeUUID(domain.UUID)
		if err!=nil {
			return nil, err
		}
		domainUUIDs[uuid] = domain.Name
	}
	return domainUUIDs, nil
}

// Apply applies the provided domain configuration to the host.
func (l *Adapter) Apply(ctx context.Context, id string, domainCfg structure.Domain) error {
	err := l.initClient()
	if err!=nil {
		return err
	}

	err = l.generator.Attach(&domainCfg)
	if err!=nil {
		return err
	}

	domain, err := l.generator.Generate(id, &domainCfg)
	if err!=nil {
		return err
	}

	domainXML, err := xml.Marshal(domain)
	if err!=nil {
		return fmt.Errorf("failed to parse generated domain xml")
	}

	_, err = l.client.DomainDefineXMLFlags(string(domainXML), libvirt.DomainDefineValidate)
	if err!=nil {
		return err
	}
	
	hotplugger := hotplug.NewLibvirtHotplugger(l.client)
	err = hotplugger.Hotplug(domain)
	if err!=nil {
		return err
	}

	return nil
}

// Destroy removes (undefines) the domain from the host. Domain must be in shutdown state for this action.
func (l *Adapter) Destroy(ctx context.Context, id string, domainCfg structure.Domain) error {
	err := l.initClient()
	if err!=nil {
		return err
	}

	uuid, err := l.parseUUID(id)
	if err!=nil {
		return err
	}
	
	domain, err := l.client.DomainLookupByUUID(uuid)
	if err!=nil {
		return err
	}

	err = l.client.DomainUndefine(domain)
	if err!=nil {
		return err
	}

	err = l.generator.Detach(&domainCfg)
	if err!=nil {
		return err
	}

	return nil
}

// Start starts the specified domain (must be defined).
func (l *Adapter) Start(ctx context.Context, id string) error {
	err := l.initClient()
	if err!=nil {
		return err
	}

	uuid, err := l.parseUUID(id)
	if err!=nil {
		return err
	}
	
	domain, err := l.client.DomainLookupByUUID(uuid)
	if err!=nil {
		return err
	}

	state, _, err := l.client.DomainGetState(domain, 0)
	if err!=nil {
		return err
	}

	if state == int32(libvirt.DomainPaused) {
		err = l.client.DomainResume(domain)
		if err!=nil {
			return err
		}
	} else {
		err = l.client.DomainCreate(domain)
		if err!=nil {
			return err
		}
	}

	return nil
}

// Reboot reboots the specified domain with the default reboot method (must be running).
func (l *Adapter) Reboot(ctx context.Context, id string) error {
	err := l.initClient()
	if err!=nil {
		return err
	}

	uuid, err := l.parseUUID(id)
	if err!=nil {
		return err
	}
	
	domain, err := l.client.DomainLookupByUUID(uuid)
	if err!=nil {
		return err
	}
	
	err = l.client.DomainReboot(domain, libvirt.DomainRebootDefault)
	if err!=nil {
		return err
	}

	return nil
}

// Pause freezes the specified domain (must be running).
func (l *Adapter) Pause(ctx context.Context, id string) error {
	err := l.initClient()
	if err!=nil {
		return err
	}

	uuid, err := l.parseUUID(id)
	if err!=nil {
		return err
	}
	
	domain, err := l.client.DomainLookupByUUID(uuid)
	if err!=nil {
		return err
	}
	
	err = l.client.DomainSuspend(domain)
	if err!=nil {
		return err
	}

	return nil
}

// Shutdown gracefully stops the domain with the default shutdown method (must be running).
func (l *Adapter) Shutdown(ctx context.Context, id string) error {
	err := l.initClient()
	if err!=nil {
		return err
	}

	uuid, err := l.parseUUID(id)
	if err!=nil {
		return err
	}
	
	domain, err := l.client.DomainLookupByUUID(uuid)
	if err!=nil {
		return err
	}
	
	err = l.client.DomainShutdown(domain)
	if err!=nil {
		return err
	}

	return nil
}

// Kill forcefully stops the domain (must be running).
func (l *Adapter) Kill(ctx context.Context, id string) error {
	err := l.initClient()
	if err!=nil {
		return err
	}

	uuid, err := l.parseUUID(id)
	if err!=nil {
		return err
	}
	
	domain, err := l.client.DomainLookupByUUID(uuid)
	if err!=nil {
		return err
	}
	
	err = l.client.DomainDestroy(domain)
	if err!=nil {
		return err
	}

	return nil
}
