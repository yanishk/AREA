import 'package:flutter/material.dart';
import 'package:auto_size_text/auto_size_text.dart';

import "../footer_page.dart";
import "myWorkflowsData.dart";
import "../globalData.dart";
import 'workflowInfos.dart';
import 'myWorkflowsFunctions.dart';

class DisplayWorkflows extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      theme: new ThemeData(
          scaffoldBackgroundColor: const Color.fromARGB(255, 255, 255, 255)),
      home: DisplayWorkflowsStateful(),
    );
  }
}

class DisplayWorkflowsStateful extends StatefulWidget {
  @override
  DisplayWorkflowsMain createState() => DisplayWorkflowsMain();
}

class DisplayWorkflowsMain extends State<DisplayWorkflowsStateful> {
  @override
  void initState() {
    super.initState();
    getWorkflows(workflowsCardList);
  }

  String apiRoute = "workflows";
  String workflowsResponseName = "workflows";
  List<Widget> workflowsCardList = [];
  String email = Globaldata.userInfos.username;

  Future<void> redirectWorkflowInfos(int index) async {
    MyWorkflowsData.selectedWorkflowsIndex = index;
    if (MyWorkflowsData.allWorkflowsData.isEmpty) {
      return;
    }
    await Navigator.push(
      Globaldata.myContext,
      MaterialPageRoute(builder: (context) => WorkflowInfosPage(),),
    );
    setState(() {
      MyWorkflowsData.allWorkflowsData.clear();
    });
    List<Widget> tmpWorkflowsCardList = [];
    await getWorkflows(tmpWorkflowsCardList);
    setState(() {
      workflowsCardList.clear();
      workflowsCardList.addAll(tmpWorkflowsCardList);
      tmpWorkflowsCardList.clear();
    });
  }

  Future<void> getWorkflows(List<Widget> cardList) async {
    final response = await MyWorkflowsfunctions.getWorkflowsResponse(apiRoute);
    if (response == null) {
      return;
    }
    setState(() {
      MyWorkflowsfunctions.storeServiceDescription(response, cardList, workflowsResponseName, redirectWorkflowInfos, email);
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: SingleChildScrollView(
        child: Column(
          children: [
            Padding(padding: EdgeInsets.only(top: 30)),
            Center(
              child: AutoSizeText("My Workflows",
                style: TextStyle(fontSize: 30,fontWeight: FontWeight.w900 ,color: Colors.black,),
                maxLines: 1,
              ),
            ),
            Column(
              children: workflowsCardList,
            ),
          ],
        ),
      ),
      bottomNavigationBar: FooterNavigationBar(currentIndex: 0),
    );
  }
}
