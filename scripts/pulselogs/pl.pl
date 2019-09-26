#!/bin/env perl

# God thanks this is not rust

=HEAD

    Script for analysing insolar node logs for consistence

=cut

use 5.018;

use Data::Dumper;
use FindBin;
use JSON::XS;


our @F = map {LogFile->new("$FindBin::Bin/../../.artifacts/launchnet/logs/discoverynodes/$_/output.log", $_)->init} (1..5);

my $mixer = new LogMixer(@F);

$mixer->get_pulse();


package ConsensusAnalyzer;

use 5.018;





package LogMixer;

use 5.018;
use Date::Parse;

sub new {
    my ($class, @inputs) = @_;
    my $self = bless {inputs => \@inputs}, $class;
}

sub get_pulse {
    my $self = shift;
    my @ret;

    for my $source (@{ $self->{inputs} }) {
        my $starttime = 0;
        while(my $obj = $source->read_line) {
            $obj->{time} = str2time($obj->{time});
            $source->putback(), last if $starttime != 0 and $obj->{time} - 10 > $starttime;
            push @ret, $obj;
        }
    }
    return @ret;
}

sub enchance {
    my $obj = shift;
    if ($obj->{message} =~ /^Consensus started/) {

    } elsif (1) {

    } else {
        die "Unhandled consensus log record: ", Dumper($obj);
    }

    return $obj;
}



package LogFile;

use 5.018;
use JSON::XS;

sub new {
    my $class = shift;
    return bless {
        fname => $_[0],
        id    => $_[1],
    };
}

sub init {
    my $self = shift;
    my $fn = $self->{fname};
    open $self->{fd}, '<', $fn or die "Can't open $fn : $!";
    $self->read_head;
    return $self;
}

sub DESTROY {
    my $self = shift;
    close $self->{fd};
}

sub read_head {
    my $self = shift;
    my $fd = $self->{fd};
    while(<$fd>) {
        $self->putback(decode_json($_)), last if /^\{/;
        $self->{head} .= $_;
    }
}

sub putback {
    my ($self, $obj) = @_;
    $self->{headbuff} = $obj;
}

sub read_line {
    my $self = shift;
    if (exists $self->{headbuff}) {
        my $obj = $self->{headbuff};
        delete $self->{headbuff};
        return $obj if $obj->{component} eq "consensus";
    }
    my $fd = $self->{fd};
    while(<$fd>) {
        next unless /^\{/;
        my $obj = decode_json($_);
        return $obj if $obj->{component} eq "consensus";
    }
}
