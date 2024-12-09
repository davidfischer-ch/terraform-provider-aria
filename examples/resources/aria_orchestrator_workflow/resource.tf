# main.tf

resource "aria_orchestrator_category" "my_company" {
  name      = "MyCompany"
  type      = "WorkflowCategory"
  parent_id = ""
}

resource "aria_orchestrator_workflow" "notification_by_mail" {
  name        = "Notification by Mail"
  description = "Send an email when a machine is provisioned."
  category_id = aria_orchestrator_category.my_company.id
  version     = "0.0.0"

  position = {
    x = 100
    y = 50
  }

  restart_mode            = 1 # resume
  resume_from_failed_mode = 0 # default

  root_name = "item13"

  input_parameters = [
    {
      name        = "inputProperties"
      description = "TODO"
      type        = "Properties"
    }
  ]

  output_parameters = []

  /*attrib:
    - value:
        string:
          value: ''
      type: string
      name: emailSubject
    - value:
        string:
          value: ''
      type: string
      name: content
    - type: 'VRA:Host'
      name: vraHost
    - value:
        string:
          value: ''
      type: string
      name: deploymentStatus
    - value:
        string:
          value: ''
      type: string
      name: requestor
    - value:
        string:
          value: ''
      type: string
      name: datacenter
    - value:
        string:
          value: ''
      type: string
      name: codeApp
    - value:
        string:
          value: ''
      type: string
      name: env
    - value:
        string:
          value: ''
      type: string
      name: osVer
    - value:
        string:
          value: ''
      type: string
      name: domain
    - value:
        string:
          value: ''
      type: string
      name: vmFlavor
    - value:
        string:
          value: ''
      type: string
      name: vmRole
    - value:
        string:
          value: ''
      type: string
      name: vlan
    - value:
        boolean:
          value: false
      type: boolean
      name: iis
    - value:
        boolean:
          value: false
      type: boolean
      name: splunk
    - value:
        boolean:
          value: false
      type: boolean
      name: controlm
    - value:
        string:
          value: ''
      type: string
      name: dotNet48core
    - value:
        string:
          value: ''
      type: string
      name: diskTabString
    - value:
        sdk-object:
          type: ResourceElement
          href: >-
            https://vralab.ceti.etat-ge.ch:443/vco/api/catalog/System/ResourceElement/9c81ce59-6d35-4d48-a0c6-c17cd72537fd/
          id: 9c81ce59-6d35-4d48-a0c6-c17cd72537fd
      type: ResourceElement
      name: logo
    - value:
        string:
          value: ''
      type: string
      name: deploymentLink
    - value:
        sdk-object:
          type: 'REST:RESTHost'
          href: >-
            https://vralab.ceti.etat-ge.ch:443/vco/api/catalog/REST/RESTHost/b303da49-127e-4f60-8d50-3898818750d4/
          id: b303da49-127e-4f60-8d50-3898818750d4
      type: 'REST:RESTHost'
      name: restServer
    - type: Array/string
      name: recipientsEmailList
    - value:
        boolean:
          value: true
      type: boolean
      name: sendToIaasAdmin
    - value:
        string:
          value: ''
      type: string
      name: smtpFromAddress
    - value:
        string:
          value: ''
      type: string
      name: smtpFromName
    - value:
        string:
          value: ''
      type: string
      name: smtpHost
    - value:
        string:
          value: ''
      type: string
      name: cmpAdminAddress
    - type: Array/string
      name: mails
    - value:
        string:
          value: ''
      type: string
      name: vmNamesString
    - value:
        string:
          value: ''
      type: string
      name: vmIpsString
    - type: number
      name: countObj*/

  /*workflow-item:
    - in-binding: {}
      out-binding: {}
      position:
        'y': 50
        x: 1200
      name: item0
      type: end
      end-mode: '0'
      comparator: 0
    - display-name: Retrieve Informations
      script:
        value: >-
          System.warn(".........................................Retrieve
          Informations.............................................");

          // GET VRA HOST

          var vraHost =
          System.getModule("ch.ocsin.dfi.core.vra.host").getvRAHost();


          // GET DEPLOYMENT INFORMATIONS

          var requestInputs = inputProperties.requestInputs;

          var deploymentId = inputProperties.deploymentId;

          var deploymentStatus = inputProperties.status;


          // GET RESOURCE VM OBJECT AS PROPERTIES

          var obj =
          System.getModule("ch.ocsin.dfi.core.vra.deployment").retrieveDeploymentResourcesVm(deploymentId);

          var vmNames = new Array();

          var vmIps = new Array();

          var countObj = obj.numberOfElements ;

          if (countObj > 1) { System.warn("More VMs found in Deployment") ;}

          for each(var cont in obj.content){
              var vmName = cont.properties.resourceName;
              var vmIp = cont.properties.address;
              System.warn("VM Name is : "+vmName+" and IP is : "+vmIp);
              vmNames.push(vmName);
              vmIps.push(vmIp);
              }

          var vmNamesString = vmNames.join(" | ");

          var vmIpsString = vmIps.join(" | ");


          // GET VM PROPERTIES

          var datacenter = requestInputs.datacenter;

          if (!datacenter) { System.warn("Datacenter is not retrieve"); }

          else { System.warn("VM Datacenter is : "+datacenter); }


          var mails = requestInputs.emails;

          if (!mails) { System.warn("mails is not retrieve"); }

          else { System.warn("Mail Notifications for LAB  : "+mails); }


          var codeApp = requestInputs.app_code;

          if (!codeApp) { System.warn("CodeApp is not retrieve"); }

          else { System.warn("Code Application is : "+codeApp); }


          var env = requestInputs.env_family;

          if (!env) { System.warn("env is not retrieve"); }

          else { System.warn("Environment is : "+env); }


          var osVer = requestInputs.os_version;

          if (!osVer) { System.warn("osVer is not retrieve"); }

          else { System.warn("OS version is : "+osVer); }


          var domain = requestInputs.domain;

          if (!domain) { System.warn("domain is not retrieve"); }

          else { System.warn("Domain is : "+domain); }


          var vmFlavor = requestInputs.flavor;

          if (!vmFlavor) { System.warn("vmFlavor is not retrieve"); }

          else { System.warn("VM flavor is : "+vmFlavor); }


          var vmRole = requestInputs.srv_role;

          if (!vmRole) { System.warn("vmRole is not retrieve"); }

          else { System.warn("VM Role is : "+vmRole); }


          var vlan = requestInputs.vlan_list;

          if (!vlan) { System.warn("vlan is not retrieve"); }

          else { System.warn("VLAN is : "+vlan); }


          var iis = requestInputs.install_iis;

          if (!iis | iis === false) { System.warn("IIS is not active"); }

          else { System.warn("IIS is active"); var iisActivation = iis; }


          var splunk = requestInputs.install_splunk;

          if (!splunk | splunk === false) { System.warn("Splunk is not active");
          }

          else { System.warn("Splunk is active"); var splunkActivation = splunk;
          }


          var controlm = requestInputs.install_controlm;

          if (!controlm | controlm === false) { System.warn("Controlm is not
          active"); }

          else { System.warn("Controlm is active"); var controlmActivation =
          controlm; }


          var dotNet48core = requestInputs.install_dotnetframework;

          if (!dotNet48core) { System.warn("DOTNET is not retrieve"); }

          else { System.warn("DOTNET is : "+dotNet48core); }


          var diskTab = new Properties();

          if (requestInputs.disk_d === true) {
              System.log("Disk D is setup");
              var disk1 = "D";
              var diskCapacity1 = requestInputs.disk_d_size;
                  diskTab.put(disk1,diskCapacity1+" Go");

              }
          if (requestInputs.disk_e === true) {
              System.log("Disk E is setup");
              var disk2 = "E";
              var diskCapacity2 = requestInputs.disk_e_size;
                  diskTab.put(disk2,diskCapacity2+" Go");
              }
          if (requestInputs.disk_f === true) {
              System.log("Disk F is setup");
              var disk3 = "F";
              var diskCapacity3 = requestInputs.disk_f_size;
                  diskTab.put(disk3,diskCapacity3+" Go");
              }
          if (requestInputs.disk_g === true) {
              System.log("Disk G is setup");
              var disk4 = "G";
              var diskCapacity4 = requestInputs.disk_g_size;
                  diskTab.put(disk4,diskCapacity4+" Go");
              }
          String

          if (!diskTab) { System.warn("Disk Tab is not retrieve"); }

          else {
              diskTabString = JSON.stringify(diskTab);
              System.warn("Disk Tab is : "+diskTabString);
              }
          // GET VRA HOST

          var vrAHost =
          System.getModule("ch.ocsin.dfi.core.vra.host").getvRAHost();


          // CREATE DEPLOYMENT LINK

          var platformUri = vrAHost.vraHost;

          var deploymentLink =
          platformUri+"/automation/#/service/catalog/consume/deployment/"+deploymentId;
        encoded: false
      in-binding:
        bind:
          - name: inputProperties
            type: Properties
            export-name: inputProperties
      out-binding:
        bind:
          - name: deploymentStatus
            type: string
            export-name: deploymentStatus
          - name: vraHost
            type: 'VRA:Host'
            export-name: vraHost
          - name: vmRole
            type: string
            export-name: vmRole
          - name: content
            type: string
            export-name: content
          - name: datacenter
            type: string
            export-name: datacenter
          - name: domain
            type: string
            export-name: domain
          - name: codeApp
            type: string
            export-name: codeApp
          - name: env
            type: string
            export-name: env
          - name: controlm
            type: boolean
            export-name: controlm
          - name: dotNet48core
            type: string
            export-name: dotNet48core
          - name: iis
            type: boolean
            export-name: iis
          - name: osVer
            type: string
            export-name: osVer
          - name: vmFlavor
            type: string
            export-name: vmFlavor
          - name: vlan
            type: string
            export-name: vlan
          - name: splunk
            type: boolean
            export-name: splunk
          - name: diskTabString
            type: string
            export-name: diskTabString
          - name: deploymentLink
            type: string
            export-name: deploymentLink
          - name: mails
            type: Array/string
            export-name: mails
          - name: vmNamesString
            type: string
            export-name: vmNamesString
          - name: vmIpsString
            type: string
            export-name: vmIpsString
          - name: countObj
            type: number
            export-name: countObj
      description: Simple task with custom script capability.
      position:
        'y': 60
        x: 340
      name: item3
      out-name: item17
      type: task
      comparator: 0
    - display-name: START w/debug
      script:
        value: |-
          try {
              var wfID      = this["workflow"].currentWorkflow.id;
              var wfName    = this["workflow"].currentWorkflow.name;
              var wfVersion = this["workflow"].currentWorkflow.version;
              var wfPath    = this["workflow"].currentWorkflow.workflowCategory.path;
              System.debug('[workflow-start] Starting workflow "' + wfPath + '/' + wfName + '" with ID:' + wfID + ', [version: ' + wfVersion + "]");
              }
          catch(e) {
              System.debug("[workflow-start] Error getting current workflow");
              }
        encoded: false
      in-binding: {}
      out-binding: {}
      description: Simple task with custom script capability.
      position:
        'y': 60
        x: 200
      name: item13
      out-name: item3
      type: task
      comparator: 0
    - display-name: Get Requestor & Build Mail Content
      script:
        value: >
          // Get Requestor

          var requestor =
          System.getContext().getParameter("__metadata_userName");

          System.log(requestor);


          // Create content Mail

          // Image du mail

          var imgMimeAttachement = logo.getContentAsMimeAttachment();


          // Gestion du titre des mails

          if ((osVer.split("-")[0]) == "Windows" ){
              if (countObj > 1 ){
                  var typeOs = "Création de Cluster Windows";
                  }
              else {
                  var typeOs = "Création de VM Windows";
                  }
              }
          else {
              if (countObj > 1 ){
                  var typeOs = "Création de Cluster Linux";
                  }
              else {
                  var typeOs = "Création de VM Linux";
                  }
              var iis = "N/A";
              var dotNet48core = "N/A";
              var dotNet48gui = "N/A";
              var diskTabString = "N/A";
              var domain = "N/A";
              }

          // les Cas du Deployments status

          if(deploymentStatus == "FINISHED"){
              var color = "style='color:rgb(39,174,96);'"
              var deploymentStatus = "Succès"
              var content = '<!DOCTYPE html><html lang="fr"><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width, initial-scale=1.0"><title>Notifications CMP</title></head><body style="font-family: Arial, sans-serif; background-color: #f3f3f3; color: #333; margin: 0; padding: 0; width: 100%;"><table width="100%" cellpadding="0" cellspacing="0" border="0" style="background-color: #f3f3f3; padding: 20px 0;"><tr><td align="center"><table width="600" cellpadding="0" cellspacing="0" border="0" style="background-color: #ffffff; border-radius: 8px; overflow: hidden; box-shadow: 0px 4px 8px rgba(0, 0, 0, 0.1);"><tr><td style="background-color: #6c757d; color: #ffffff; text-align: center; padding: 20px;"><h1 style="font-size: 1.5rem; font-weight: normal; margin: 0;">'+typeOs+'</h1></td></tr><tr><td style="padding: 20px;"><h4 style="text-align: center; color: #555; font-weight: bold; margin: 0 0 20px; font-size: 1.2rem;">Les caractéristiques de la VM :</h4><table width="100%" cellpadding="0" cellspacing="0" border="0" style="margin: 0 auto; border-collapse: collapse;"><tr><th style="padding: 15px; background-color: #e9ecef; color: #333; text-align: left; border-bottom: 1px solid #e0e0e0; font-weight: bold;">Nom de VM</th><th style="padding: 15px; background-color: #e9ecef; color: #333; text-align: left; border-bottom: 1px solid #e0e0e0; font-weight: bold;">Deploiement Link</th><th style="padding: 15px; background-color: #e9ecef; color: #333; text-align: left; border-bottom: 1px solid #e0e0e0; font-weight: bold;">Adresse IP</th><th style="padding: 15px; background-color: #e9ecef; color: #333; text-align: left; border-bottom: 1px solid #e0e0e0; font-weight: bold;">Status du Déploiement</th></tr><tr><td style="padding: 15px; background-color: #ffffff; border-bottom: 1px solid #e0e0e0;">'+vmName+'</td><td style="padding: 15px; background-color: #ffffff; border-bottom: 1px solid #e0e0e0;"><a href='+deploymentLink+'>Link</a></td><td style="padding: 15px; background-color: #ffffff; border-bottom: 1px solid #e0e0e0;">'+vmIP+'</td><td style="padding: 15px; background-color: #ffffff; border-bottom: 1px solid #e0e0e0;"><strong '+color+'>'+deploymentStatus+'</strong></td></tr></table></td></tr><tr><td style="padding: 20px;"><h4 style="text-align: center; color: #555; font-weight: bold; margin: 0 0 20px; font-size: 1.2rem;">Les détails de la Demande :</h4><table width="100%" cellpadding="0" cellspacing="0" border="0" style="margin: 0 auto; border-collapse: collapse;"><tr><th style="padding: 15px; background-color: #e9ecef; color: #333; text-align: left; border-bottom: 1px solid #e0e0e0; font-weight: bold;">Code Applicatif</th><th style="padding: 15px; background-color: #e9ecef; color: #333; text-align: left; border-bottom: 1px solid #e0e0e0; font-weight: bold;">Environnement</th><th style="padding: 15px; background-color: #e9ecef; color: #333; text-align: left; border-bottom: 1px solid #e0e0e0; font-weight: bold;">Version OS</th><th style="padding: 15px; background-color: #e9ecef; color: #333; text-align: left; border-bottom: 1px solid #e0e0e0; font-weight: bold;">Domain</th><th style="padding: 15px; background-color: #e9ecef; color: #333; text-align: left; border-bottom: 1px solid #e0e0e0; font-weight: bold;">Taille</th><th style="padding: 15px; background-color: #e9ecef; color: #333; text-align: left; border-bottom: 1px solid #e0e0e0; font-weight: bold;">Rôle</th><th style="padding: 15px; background-color: #e9ecef; color: #333; text-align: left; border-bottom: 1px solid #e0e0e0; font-weight: bold;">vLAN</th><th style="padding: 15px; background-color: #e9ecef; color: #333; text-align: left; border-bottom: 1px solid #e0e0e0; font-weight: bold;">Datacenter</th><th style="padding: 15px; background-color: #e9ecef; color: #333; text-align: left; border-bottom: 1px solid #e0e0e0; font-weight: bold;">IIS</th><th style="padding: 15px; background-color: #e9ecef; color: #333; text-align: left; border-bottom: 1px solid #e0e0e0; font-weight: bold;">Splunk</th><th style="padding: 15px; background-color: #e9ecef; color: #333; text-align: left; border-bottom: 1px solid #e0e0e0; font-weight: bold;">Control-M</th><th style="padding: 15px; background-color: #e9ecef; color: #333; text-align: left; border-bottom: 1px solid #e0e0e0; font-weight: bold;">.NET Framework</th><th style="padding: 15px; background-color: #e9ecef; color: #333; text-align: left; border-bottom: 1px solid #e0e0e0; font-weight: bold;">Disques</th></tr><tr><td style="padding: 15px; background-color: #ffffff; border-bottom: 1px solid #e0e0e0;">'+codeApp+'</td><td style="padding: 15px; background-color: #ffffff; border-bottom: 1px solid #e0e0e0;">'+env+'</td><td style="padding: 15px; background-color: #ffffff; border-bottom: 1px solid #e0e0e0;">'+osVer+'</td><td style="padding: 15px; background-color: #ffffff; border-bottom: 1px solid #e0e0e0;">'+domain+'</td><td style="padding: 15px; background-color: #ffffff; border-bottom: 1px solid #e0e0e0;">'+vmFlavor+'</td><td style="padding: 15px; background-color: #ffffff; border-bottom: 1px solid #e0e0e0;">'+vmRole+'</td><td style="padding: 15px; background-color: #ffffff; border-bottom: 1px solid #e0e0e0;">'+vlan+'</td><td style="padding: 15px; background-color: #ffffff; border-bottom: 1px solid #e0e0e0;">'+datacenter+'</td><td style="padding: 15px; background-color: #ffffff; border-bottom: 1px solid #e0e0e0;">'+iis+'</td><td style="padding: 15px; background-color: #ffffff; border-bottom: 1px solid #e0e0e0;">'+splunk+'</td><td style="padding: 15px; background-color: #ffffff; border-bottom: 1px solid #e0e0e0;">'+controlm+'</td><td style="padding: 15px; background-color: #ffffff; border-bottom: 1px solid #e0e0e0;">'+dotNet48core+'</td><td style="padding: 15px; background-color: #ffffff; border-bottom: 1px solid #e0e0e0;">'+diskTabString+'</td></tr></table></td></tr><tr><td style="background-color: #333333; color: #ffffff; text-align: center; padding: 20px; font-size: 0.9rem;">En cas de besoin, merci de contacter le support CMP : <a href="mailto:secteur.csc@etat.ge.ch" style="color: #007bff; text-decoration: none;">secteur.csc@etat.ge.ch</a></td></tr></table></td></tr></table></body></html>';
           if (countObj > 1 ){
               // Create Subject
               emailSubject = "[CMP] Votre Cluster "+vmName+" est maintenant disponible ";
              }
           else {
                    // Create Subject
               emailSubject = "[CMP] Votre VM "+vmName+" est maintenant disponible ";
               }
          }

          else if (deploymentStatus == "CANCELLED"){
              var vmName = "N/A"
              var vmIP = "N/A"
              var color = "style='color:rgb(240,178,122);'"
              var deploymentStatus = "Annulé";
              var content = '<!DOCTYPE html><html lang="fr"><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width, initial-scale=1.0"><title>Notifications CMP</title></head><body style="font-family: Arial, sans-serif; background-color: #f3f3f3; color: #333; margin: 0; padding: 0; width: 100%;"><table width="100%" cellpadding="0" cellspacing="0" border="0" style="background-color: #f3f3f3; padding: 20px 0;"><tr><td align="center"><table width="600" cellpadding="0" cellspacing="0" border="0" style="background-color: #ffffff; border-radius: 8px; overflow: hidden; box-shadow: 0px 4px 8px rgba(0, 0, 0, 0.1);"><tr><td style="background-color: #6c757d; color: #ffffff; text-align: center; padding: 20px;"><h1 style="font-size: 1.5rem; font-weight: normal; margin: 0;">'+typeOs+'</h1></td></tr><tr><td style="padding: 20px;"><h4 style="text-align: center; color: #555; font-weight: bold; margin: 0 0 20px; font-size: 1.2rem;">Les caractéristiques de la VM :</h4><table width="100%" cellpadding="0" cellspacing="0" border="0" style="margin: 0 auto; border-collapse: collapse;"><tr><th style="padding: 15px; background-color: #e9ecef; color: #333; text-align: left; border-bottom: 1px solid #e0e0e0; font-weight: bold;">Nom de VM</th><th style="padding: 15px; background-color: #e9ecef; color: #333; text-align: left; border-bottom: 1px solid #e0e0e0; font-weight: bold;">Deploiement Link</th><th style="padding: 15px; background-color: #e9ecef; color: #333; text-align: left; border-bottom: 1px solid #e0e0e0; font-weight: bold;">Adresse IP</th><th style="padding: 15px; background-color: #e9ecef; color: #333; text-align: left; border-bottom: 1px solid #e0e0e0; font-weight: bold;">Status du Déploiement</th></tr><tr><td style="padding: 15px; background-color: #ffffff; border-bottom: 1px solid #e0e0e0;">'+vmName+'</td><td style="padding: 15px; background-color: #ffffff; border-bottom: 1px solid #e0e0e0;"><a href='+deploymentLink+'>Link</a></td><td style="padding: 15px; background-color: #ffffff; border-bottom: 1px solid #e0e0e0;">'+vmIP+'</td><td style="padding: 15px; background-color: #ffffff; border-bottom: 1px solid #e0e0e0;"><strong '+color+'>'+deploymentStatus+'</strong></td></tr></table></td></tr><tr><td style="padding: 20px;"><h4 style="text-align: center; color: #555; font-weight: bold; margin: 0 0 20px; font-size: 1.2rem;">Les détails de la Demande :</h4><table width="100%" cellpadding="0" cellspacing="0" border="0" style="margin: 0 auto; border-collapse: collapse;"><tr><th style="padding: 15px; background-color: #e9ecef; color: #333; text-align: left; border-bottom: 1px solid #e0e0e0; font-weight: bold;">Code Applicatif</th><th style="padding: 15px; background-color: #e9ecef; color: #333; text-align: left; border-bottom: 1px solid #e0e0e0; font-weight: bold;">Environnement</th><th style="padding: 15px; background-color: #e9ecef; color: #333; text-align: left; border-bottom: 1px solid #e0e0e0; font-weight: bold;">Version OS</th><th style="padding: 15px; background-color: #e9ecef; color: #333; text-align: left; border-bottom: 1px solid #e0e0e0; font-weight: bold;">Domain</th><th style="padding: 15px; background-color: #e9ecef; color: #333; text-align: left; border-bottom: 1px solid #e0e0e0; font-weight: bold;">Taille</th><th style="padding: 15px; background-color: #e9ecef; color: #333; text-align: left; border-bottom: 1px solid #e0e0e0; font-weight: bold;">Rôle</th><th style="padding: 15px; background-color: #e9ecef; color: #333; text-align: left; border-bottom: 1px solid #e0e0e0; font-weight: bold;">vLAN</th><th style="padding: 15px; background-color: #e9ecef; color: #333; text-align: left; border-bottom: 1px solid #e0e0e0; font-weight: bold;">Datacenter</th><th style="padding: 15px; background-color: #e9ecef; color: #333; text-align: left; border-bottom: 1px solid #e0e0e0; font-weight: bold;">IIS</th><th style="padding: 15px; background-color: #e9ecef; color: #333; text-align: left; border-bottom: 1px solid #e0e0e0; font-weight: bold;">Splunk</th><th style="padding: 15px; background-color: #e9ecef; color: #333; text-align: left; border-bottom: 1px solid #e0e0e0; font-weight: bold;">Control-M</th><th style="padding: 15px; background-color: #e9ecef; color: #333; text-align: left; border-bottom: 1px solid #e0e0e0; font-weight: bold;">.NET Framework</th><th style="padding: 15px; background-color: #e9ecef; color: #333; text-align: left; border-bottom: 1px solid #e0e0e0; font-weight: bold;">Disques</th></tr><tr><td style="padding: 15px; background-color: #ffffff; border-bottom: 1px solid #e0e0e0;">'+codeApp+'</td><td style="padding: 15px; background-color: #ffffff; border-bottom: 1px solid #e0e0e0;">'+env+'</td><td style="padding: 15px; background-color: #ffffff; border-bottom: 1px solid #e0e0e0;">'+osVer+'</td><td style="padding: 15px; background-color: #ffffff; border-bottom: 1px solid #e0e0e0;">'+domain+'</td><td style="padding: 15px; background-color: #ffffff; border-bottom: 1px solid #e0e0e0;">'+vmFlavor+'</td><td style="padding: 15px; background-color: #ffffff; border-bottom: 1px solid #e0e0e0;">'+vmRole+'</td><td style="padding: 15px; background-color: #ffffff; border-bottom: 1px solid #e0e0e0;">'+vlan+'</td><td style="padding: 15px; background-color: #ffffff; border-bottom: 1px solid #e0e0e0;">'+datacenter+'</td><td style="padding: 15px; background-color: #ffffff; border-bottom: 1px solid #e0e0e0;">'+iis+'</td><td style="padding: 15px; background-color: #ffffff; border-bottom: 1px solid #e0e0e0;">'+splunk+'</td><td style="padding: 15px; background-color: #ffffff; border-bottom: 1px solid #e0e0e0;">'+controlm+'</td><td style="padding: 15px; background-color: #ffffff; border-bottom: 1px solid #e0e0e0;">'+dotNet48core+'</td><td style="padding: 15px; background-color: #ffffff; border-bottom: 1px solid #e0e0e0;">'+diskTabString+'</td></tr></table></td></tr><tr><td style="background-color: #333333; color: #ffffff; text-align: center; padding: 20px; font-size: 0.9rem;">En cas de besoin, merci de contacter le support CMP : <a href="mailto:secteur.csc@etat.ge.ch" style="color: #007bff; text-decoration: none;">secteur.csc@etat.ge.ch</a></td></tr></table></td></tr></table></body></html>';
          // Create Subject
              emailSubject = "[CMP] Votre déploiement a été annulé " ;
          }else {
              var deploymentStatus = "Echec"
              var color = "style='color:rgb(231,76,60);'"
              var content = '<!DOCTYPE html><html lang="fr"><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width, initial-scale=1.0"><title>Notifications CMP</title></head><body style="font-family: Arial, sans-serif; background-color: #f3f3f3; color: #333; margin: 0; padding: 0; width: 100%;"><table width="100%" cellpadding="0" cellspacing="0" border="0" style="background-color: #f3f3f3; padding: 20px 0;"><tr><td align="center"><table width="600" cellpadding="0" cellspacing="0" border="0" style="background-color: #ffffff; border-radius: 8px; overflow: hidden; box-shadow: 0px 4px 8px rgba(0, 0, 0, 0.1);"><tr><td style="background-color: #6c757d; color: #ffffff; text-align: center; padding: 20px;"><h1 style="font-size: 1.5rem; font-weight: normal; margin: 0;">'+typeOs+'</h1></td></tr><tr><td style="padding: 20px;"><h4 style="text-align: center; color: #555; font-weight: bold; margin: 0 0 20px; font-size: 1.2rem;">Les caractéristiques de la VM :</h4><table width="100%" cellpadding="0" cellspacing="0" border="0" style="margin: 0 auto; border-collapse: collapse;"><tr><th style="padding: 15px; background-color: #e9ecef; color: #333; text-align: left; border-bottom: 1px solid #e0e0e0; font-weight: bold;">Nom de VM</th><th style="padding: 15px; background-color: #e9ecef; color: #333; text-align: left; border-bottom: 1px solid #e0e0e0; font-weight: bold;">Deploiement Link</th><th style="padding: 15px; background-color: #e9ecef; color: #333; text-align: left; border-bottom: 1px solid #e0e0e0; font-weight: bold;">Adresse IP</th><th style="padding: 15px; background-color: #e9ecef; color: #333; text-align: left; border-bottom: 1px solid #e0e0e0; font-weight: bold;">Status du Déploiement</th></tr><tr><td style="padding: 15px; background-color: #ffffff; border-bottom: 1px solid #e0e0e0;">'+vmName+'</td><td style="padding: 15px; background-color: #ffffff; border-bottom: 1px solid #e0e0e0;"><a href='+deploymentLink+'>Link</a></td><td style="padding: 15px; background-color: #ffffff; border-bottom: 1px solid #e0e0e0;">'+vmIP+'</td><td style="padding: 15px; background-color: #ffffff; border-bottom: 1px solid #e0e0e0;"><strong '+color+'>'+deploymentStatus+'</strong></td></tr></table></td></tr><tr><td style="padding: 20px;"><h4 style="text-align: center; color: #555; font-weight: bold; margin: 0 0 20px; font-size: 1.2rem;">Les détails de la Demande :</h4><table width="100%" cellpadding="0" cellspacing="0" border="0" style="margin: 0 auto; border-collapse: collapse;"><tr><th style="padding: 15px; background-color: #e9ecef; color: #333; text-align: left; border-bottom: 1px solid #e0e0e0; font-weight: bold;">Code Applicatif</th><th style="padding: 15px; background-color: #e9ecef; color: #333; text-align: left; border-bottom: 1px solid #e0e0e0; font-weight: bold;">Environnement</th><th style="padding: 15px; background-color: #e9ecef; color: #333; text-align: left; border-bottom: 1px solid #e0e0e0; font-weight: bold;">Version OS</th><th style="padding: 15px; background-color: #e9ecef; color: #333; text-align: left; border-bottom: 1px solid #e0e0e0; font-weight: bold;">Domain</th><th style="padding: 15px; background-color: #e9ecef; color: #333; text-align: left; border-bottom: 1px solid #e0e0e0; font-weight: bold;">Taille</th><th style="padding: 15px; background-color: #e9ecef; color: #333; text-align: left; border-bottom: 1px solid #e0e0e0; font-weight: bold;">Rôle</th><th style="padding: 15px; background-color: #e9ecef; color: #333; text-align: left; border-bottom: 1px solid #e0e0e0; font-weight: bold;">vLAN</th><th style="padding: 15px; background-color: #e9ecef; color: #333; text-align: left; border-bottom: 1px solid #e0e0e0; font-weight: bold;">Datacenter</th><th style="padding: 15px; background-color: #e9ecef; color: #333; text-align: left; border-bottom: 1px solid #e0e0e0; font-weight: bold;">IIS</th><th style="padding: 15px; background-color: #e9ecef; color: #333; text-align: left; border-bottom: 1px solid #e0e0e0; font-weight: bold;">Splunk</th><th style="padding: 15px; background-color: #e9ecef; color: #333; text-align: left; border-bottom: 1px solid #e0e0e0; font-weight: bold;">Control-M</th><th style="padding: 15px; background-color: #e9ecef; color: #333; text-align: left; border-bottom: 1px solid #e0e0e0; font-weight: bold;">.NET Framework</th><th style="padding: 15px; background-color: #e9ecef; color: #333; text-align: left; border-bottom: 1px solid #e0e0e0; font-weight: bold;">Disques</th></tr><tr><td style="padding: 15px; background-color: #ffffff; border-bottom: 1px solid #e0e0e0;">'+codeApp+'</td><td style="padding: 15px; background-color: #ffffff; border-bottom: 1px solid #e0e0e0;">'+env+'</td><td style="padding: 15px; background-color: #ffffff; border-bottom: 1px solid #e0e0e0;">'+osVer+'</td><td style="padding: 15px; background-color: #ffffff; border-bottom: 1px solid #e0e0e0;">'+domain+'</td><td style="padding: 15px; background-color: #ffffff; border-bottom: 1px solid #e0e0e0;">'+vmFlavor+'</td><td style="padding: 15px; background-color: #ffffff; border-bottom: 1px solid #e0e0e0;">'+vmRole+'</td><td style="padding: 15px; background-color: #ffffff; border-bottom: 1px solid #e0e0e0;">'+vlan+'</td><td style="padding: 15px; background-color: #ffffff; border-bottom: 1px solid #e0e0e0;">'+datacenter+'</td><td style="padding: 15px; background-color: #ffffff; border-bottom: 1px solid #e0e0e0;">'+iis+'</td><td style="padding: 15px; background-color: #ffffff; border-bottom: 1px solid #e0e0e0;">'+splunk+'</td><td style="padding: 15px; background-color: #ffffff; border-bottom: 1px solid #e0e0e0;">'+controlm+'</td><td style="padding: 15px; background-color: #ffffff; border-bottom: 1px solid #e0e0e0;">'+dotNet48core+'</td><td style="padding: 15px; background-color: #ffffff; border-bottom: 1px solid #e0e0e0;">'+diskTabString+'</td></tr></table></td></tr><tr><td style="background-color: #333333; color: #ffffff; text-align: center; padding: 20px; font-size: 0.9rem;">En cas de besoin, merci de contacter le support CMP : <a href="mailto:secteur.csc@etat.ge.ch" style="color: #007bff; text-decoration: none;">secteur.csc@etat.ge.ch</a></td></tr></table></td></tr></table></body></html>';
           // Create Subject
              emailSubject = "[CMP] Votre déploiement est en echec " ;
              }
          System.debug(content);


          var mailRequestor =
          System.getModule("ch.ocsin.dfi.core.vro.mail").getUserMail(restServer,
          requestor);

          System.warn("Requestor mail is "+mailRequestor);


          var recipientsEmailList = new Array();


          if (!mails) {
              System.log("This deployment is not for LAB env");
              }
          else {
              if (mails.length > 0) {
                  for each(var mail in mails){
                      System.log(mail);
                      recipientsEmailList.push(mail);
                      }
                  }
              }

          // GET ALL MAILS IN ONE TAB

          recipientsEmailList.push(mailRequestor);

          //recipientsEmailList.push("gedeon.adechi-bobo@etat.ge.ch");

          //recipientsEmailList.push(cmpAdminAddress);
        encoded: false
      in-binding:
        bind:
          - name: vmFlavor
            type: string
            export-name: vmFlavor
          - name: codeApp
            type: string
            export-name: codeApp
          - name: vmName
            type: string
            export-name: vmNamesString
          - name: vmRole
            type: string
            export-name: vmRole
          - name: datacenter
            type: string
            export-name: datacenter
          - name: domain
            type: string
            export-name: domain
          - name: vmIP
            type: string
            export-name: vmIpsString
          - name: vlan
            type: string
            export-name: vlan
          - name: splunk
            type: boolean
            export-name: splunk
          - name: iis
            type: boolean
            export-name: iis
          - name: env
            type: string
            export-name: env
          - name: dotNet48core
            type: string
            export-name: dotNet48core
          - name: controlm
            type: boolean
            export-name: controlm
          - name: osVer
            type: string
            export-name: osVer
          - name: deploymentStatus
            type: string
            export-name: deploymentStatus
          - name: diskTabString
            type: string
            export-name: diskTabString
          - name: logo
            type: ResourceElement
            export-name: logo
          - name: deploymentLink
            type: string
            export-name: deploymentLink
          - name: sendToIaasAdmin
            type: boolean
            export-name: sendToIaasAdmin
          - name: cmpAdminAddress
            type: string
            export-name: cmpAdminAddress
          - name: restServer
            type: 'REST:RESTHost'
            export-name: restServer
          - name: mails
            type: Array/string
            export-name: mails
          - name: countObj
            type: number
            export-name: countObj
      out-binding:
        bind:
          - name: content
            type: string
            export-name: content
          - name: requestor
            type: string
            export-name: requestor
          - name: recipientsEmailList
            type: Array/string
            export-name: recipientsEmailList
          - name: emailSubject
            type: string
            export-name: emailSubject
      description: Simple task with custom script capability.
      position:
        'y': 60
        x: 720
      name: item14
      out-name: item15
      type: task
      comparator: 0
    - display-name: Send Mail
      runtime: 'powercli:13-powershell-7.4'
      script:
        value: |
          function Handler($context, $inputs) {
              # Get Mail Inputs
              $content = $inputs.content
              $smtpFromAddress = $inputs.smtpFromAddress
              $emailSubject = $inputs.emailSubject
              $recipientsEmailList = $inputs.recipientsEmailList
              $smtpHost = $inputs.smtpHost
              # Get user
              #$userMail = get-aduser -Identity "$requestor".replace("vsr","") -Property Emailaddress).EmailAddress
              # Send Mail
              foreach ($mail in $recipientsEmailList){
              $sendMail = send-MailMessage -From $smtpFromAddress -to $mail -Subject $emailSubject -Body $content -BodyAsHtml -SmtpServer $smtpHost -encoding utf8
              }
          }
        encoded: false
      in-binding:
        bind:
          - name: content
            type: string
            export-name: content
          - name: smtpFromAddress
            type: string
            export-name: smtpFromAddress
          - name: smtpFromName
            type: string
            export-name: smtpFromName
          - name: emailSubject
            type: string
            export-name: emailSubject
          - name: recipientsEmailList
            type: Array/string
            export-name: recipientsEmailList
          - name: smtpHost
            type: string
            export-name: smtpHost
      out-binding: {}
      description: Simple task with custom script capability.
      position:
        'y': 60
        x: 960
      name: item15
      out-name: item0
      type: task
      comparator: 0
    - display-name: Retrieve SMTP Parameters From vRO Configurations
      script:
        value: >-
          System.warn("Retrieve SMTP Properties");

          var emailConfigPath = "/OCSIN/Email";

          var emailConfigName ="Email Server Configuration";

          var emailAttributes =
          System.getModule("ch.ocsin.dfi.core.vro.helpers").getAttributesValuesFromPathAndNameAsProp(emailConfigPath,emailConfigName);


          smtpFromAddress = emailAttributes.get("cmpNotificationAddress");

          System.warn("   - From Address: " + smtpFromAddress);


          smtpFromName = emailAttributes.get("cmpNotificationName");

          System.warn("   - From Name: " + smtpFromName);


          smtpHost = emailAttributes.get("smtpHost");

          System.warn("   - SMTP Host: " + smtpHost);

          cmpAdminAddress = emailAttributes.get("cmpAdminAddress");

          System.warn("   - Castle Admin Address: " + cmpAdminAddress);
        encoded: false
      in-binding: {}
      out-binding:
        bind:
          - name: smtpFromAddress
            type: string
            export-name: smtpFromAddress
          - name: smtpFromName
            type: string
            export-name: smtpFromName
          - name: smtpHost
            type: string
            export-name: smtpHost
          - name: cmpAdminAddress
            type: string
            export-name: cmpAdminAddress
      description: Simple task with custom script capability.
      position:
        'y': 60
        x: 500
      name: item17
      out-name: item14
      type: task
      comparator: 0*/

  /*
  inputForms:
    - layout:
        pages:
          - id: page_1
            title: General
            sections:
              - id: section_0
                fields:
                  - id: inputProperties
                    display: datagrid
      schema:
        inputProperties:
          type:
            dataType: complex
            isMultiple: true
            fields:
              - type:
                  dataType: string
                  isMultiple: false
                label: Key
                id: key
                constraints: {}
              - type:
                  dataType: string
                  isMultiple: false
                label: Value
                id: value
                constraints: {}
          label: inputProperties
          id: inputProperties
          constraints: {}
      itemId: ''
  */
}
