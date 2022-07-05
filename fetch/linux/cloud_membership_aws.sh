#!/bin/bash

# Copyright (c) 2021 Sorint.lab S.p.A.
#
# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.
#
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU General Public License for more details.
#
# You should have received a copy of the GNU General Public License
# along with this program.  If not, see <http://www.gnu.org/licenses/>.

# https://serverfault.com/questions/462903/how-to-know-if-a-machine-is-an-ec2-instance

# This first, simple check will work for many older instance types.
if [ -f /sys/hypervisor/uuid ]; then
  # File should be readable by non-root users.
  if [ "$(head -c 3 /sys/hypervisor/uuid)" = "ec2" ]; then
      echo "true"
  else
      echo "false"
  fi

# This check will work on newer m5/c5 instances, but only if you have root!
elif [ -r /sys/devices/virtual/dmi/id/product_uuid ]; then
  # If the file exists AND is readable by us, we can rely on it.
  if [ "$(head -c 3 /sys/devices/virtual/dmi/id/product_uuid)" = "EC2" ]; then
      echo "true"
  else
      echo "false"
  fi

else
  # Fallback check of http://169.254.169.254/. If we wanted to be REALLY
  # authoritative, we could follow Amazon's suggestions for cryptographically
  # verifying their signature, see here:
  #    https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/instance-identity-documents.html
  # but this is almost certainly overkill for this purpose (and the above
  # checks of "EC2" prefixes have a higher false positive potential, anyway).
  if "$(curl -s -m 5 http://169.254.169.254/latest/dynamic/instance-identity/document | grep -q availabilityZone)" ; then
      echo "true"
  else
      echo "false"
  fi

fi
