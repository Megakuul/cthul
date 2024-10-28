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

// Start starts the specified domain (must be defined).
func (l *LibvirtController) Start(id string) error {
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
func (l *LibvirtController) Reboot(id string) error {
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
func (l *LibvirtController) Pause(id string) error {
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
func (l *LibvirtController) Resume(id string) error {
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
func (l *LibvirtController) Shutdown(id string) error {
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
func (l *LibvirtController) Kill(id string) error {
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
