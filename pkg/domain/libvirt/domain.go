/**
 * Cthul System
 *
 * Copyright (C) 2024 Linus Ilian Moser <linus.moser@megakuul.ch>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package libvirt

import (
	"context"
	"encoding/xml"
	"fmt"

	"cthul.io/cthul/pkg/domain/structure"
	"github.com/digitalocean/go-libvirt"
)

// ListDomains fetches the uuid of all domains located on this node.
func (l *LibvirtController) ListDomains(ctx context.Context) ([]string, error) {
	err := l.initClient()
	if err!=nil {
		return nil, err
	}

	domains, _, err := l.client.ConnectListAllDomains(-1, 0)
	if err!=nil {
		return nil, err
	}

	domainUUIDs := []string{}
	for _, domain := range domains {
		uuid, err := l.serializeUUID(domain.UUID)
		if err!=nil {
			return nil, err
		}
		domainUUIDs = append(domainUUIDs, uuid)
	}
	return domainUUIDs, nil
}

// ApplyDomain applies the provided domain configuration to the host.
func (l *LibvirtController) ApplyDomain(ctx context.Context, domainCfg structure.Domain) error {
	err := l.initClient()
	if err!=nil {
		return err
	}

	for _, device := range domainCfg.BlockDevices {
		// PoC: l.granit.AttachBlock(device.GranitBlockDeviceId)
		_ = device
	}

	for _, device := range domainCfg.NetworkDevices {
		// PoC: l.proton.AttachInterface(device.ProtonNetworkDeviceId)
		_ = device
	}

	for _, device := range domainCfg.SerialDevices {
		// PoC: l.wave.AttachSerial(device.WaveSerialDeviceId)
		_ = device
	}

	for _, device := range domainCfg.GraphicDevices {
		// PoC: l.wave.AttachGraphic(device.WaveGraphicDeviceId)
		_ = device
	}
	err = l.generator.Prepare(domainCfg)
	if err!=nil {
		return err
	}

	domain, err := l.generator.Generate(domainCfg)
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

	err = l.hotplugger.Hotplug(domain)
	if err!=nil {
		return err
	}

	return nil
}

// DestroyDomain removes (undefines) the domain from the host. Domain must be in shutdown state for this action.
func (l *LibvirtController) DestroyDomain(ctx context.Context, domainCfg structure.Domain) error {
	err := l.initClient()
	if err!=nil {
		return err
	}

	err = l.generator.Release(domainCfg)
	if err!=nil {
		return err
	}
	
	for _, device := range domainCfg.BlockDevices {
		// PoC: l.granit.ReleaseBlock(device.GranitBlockDeviceId)
		_ = device
	}

	for _, device := range domainCfg.NetworkDevices {
		// PoC: l.proton.ReleaseInterface(device.ProtonNetworkDeviceId)
		_ = device
	}

	for _, device := range domainCfg.SerialDevices {
		// PoC: l.wave.ReleaseSerial(device.WaveSerialDeviceId)
		_ = device
	}

	for _, device := range domainCfg.GraphicDevices {
		// PoC: l.wave.ReleaseGraphic(device.WaveGraphicDeviceId)
		_ = device
	}

	uuid, err := l.parseUUID(domainCfg.UUID)
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

	return nil
}

// Start starts the specified domain (must be defined).
func (l *LibvirtController) StartDomain(ctx context.Context, id string) error {
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
	
	err = l.client.DomainCreate(domain)
	if err!=nil {
		return err
	}

	return nil
}

// Reboot reboots the specified domain with the default reboot method (must be running).
func (l *LibvirtController) RebootDomain(ctx context.Context, id string) error {
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
func (l *LibvirtController) PauseDomain(ctx context.Context, id string) error {
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

// Resume unpauses the specified domain (must be paused).
func (l *LibvirtController) ResumeDomain(ctx context.Context, id string) error {
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
	
	err = l.client.DomainResume(domain)
	if err!=nil {
		return err
	}

	return nil
}

// Shutdown gracefully stops the domain with the default shutdown method (must be running).
func (l *LibvirtController) ShutdownDomain(ctx context.Context, id string) error {
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
func (l *LibvirtController) KillDomain(ctx context.Context, id string) error {
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
